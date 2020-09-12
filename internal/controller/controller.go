package controller

import (
	context "context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"

	"hurracloud.io/jawhar/internal/agent"
	pb "hurracloud.io/jawhar/internal/agent/proto"
	"hurracloud.io/jawhar/internal/database"
	"hurracloud.io/jawhar/internal/models"
)

/* GET /sources */
func GetSources(c echo.Context) error {

	response, err := agent.Client.GetDrives(context.Background(), &pb.GetDrivesRequest{})
	if err != nil {
		log.Error("Agent Client Could not GetDrives: ", err)
	}

	for _, drive := range response.Drives {
		// Check if we know this drive
		log.Tracef("Agent returned drive %v", drive)
		var aDrive models.Drive
		result := database.DB.Where("serial_number = ?", drive.SerialNumber).First(&aDrive)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// First time we see this drive
			log.Debugf("First time we see drive with serial number: %v. Creating a record.", drive.SerialNumber)
			aDrive = models.Drive{SerialNumber: drive.SerialNumber}
			database.DB.Create(&aDrive)
		}

		aDrive.Name = drive.Name
		aDrive.DeviceFile = drive.DeviceFile
		aDrive.DriveType = drive.Type
		aDrive.SizeBytes = drive.SizeBytes
		database.DB.Save(&aDrive)

		for _, partition := range drive.Partitions {
			// Check if we know this partition
			log.Tracef("Agent returned partition %v", partition)
			var aPartition models.DrivePartition
			unique_name := fmt.Sprintf("%s-%s", drive.SerialNumber, partition.Name)
			result := database.DB.Where("name = ?", unique_name).First(&aPartition)

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// First time we see this partition, create it
				log.Debugf("First time we see partition with this name: %v. Creating a record.", unique_name)
				aPartition = models.DrivePartition{Name: unique_name}
				aPartition.Drive = aDrive
				if partition.Label != "" {
					aPartition.Caption = partition.Label
				} else {
					aPartition.Caption = partition.Name
				}
				database.DB.Create(&aPartition)
			}

			aPartition.DeviceFile = partition.DeviceFile
			aPartition.SizeBytes = partition.SizeBytes
			aPartition.AvailableBytes = partition.AvailableBytes
			aPartition.MountPoint = partition.MountPoint
			aPartition.Label = partition.Label
			aPartition.IsReadOnly = partition.IsReadOnly
			aPartition.Status = "unmounted"
			if partition.MountPoint != "" {
				aPartition.Status = "mounted"
			}
			database.DB.Save(&aPartition)
		}

	}

	var partitions []models.DrivePartition
	database.DB.Preload("Drive").Find(&partitions)
	return c.JSON(http.StatusOK, partitions)
}
