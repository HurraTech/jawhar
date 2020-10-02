package manager

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"hurracloud.io/jawhar/internal/agent"
	pb "hurracloud.io/jawhar/internal/agent/proto"
	"hurracloud.io/jawhar/internal/database"
	"hurracloud.io/jawhar/internal/models"
	"hurracloud.io/jawhar/internal/zahif"
	zahif_pb "hurracloud.io/jawhar/internal/zahif/proto"
)

var systemPartitions = map[string]bool{
	"/":         true,
	"/data":     true,
	"/boot":     true,
	"/boot/efi": true,
	"/uboot":    true,
}

var systemDevices = map[string]bool{
	// RPI
	"/dev/mmcblk0p1": true,
	"/dev/mmcblk0p2": true,
	"/dev/mmcblk0p3": true,
	"/dev/mmcblk0p4": true,

	// Thinkpad
	"/dev/nvme0n1p2": true,
	"/dev/nvme0n1p3": true,
	"/dev/nvme0n1p4": true,
}

var supportedFilesystems = map[string]bool{
	"vfat": true,
	"ext4": true,
	"ext3": true,
	"ntfs": true,
}

func UpdateSources(internalStoragePath string) error {
	if err := retrieveSources(); err != nil {
		return err
	}

	updateInternalStorageDummyPartition(internalStoragePath)
	updateIndexingProgress()
	return nil
}

func retrieveSources() error {
	response, err := agent.Client.GetDrives(context.Background(), &pb.GetDrivesRequest{})
	if err != nil {
		return fmt.Errorf("Agent Client Failed to call GetDrives: %s", err)
	}

	var attacheDrivesSN []string
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
		aDrive.Status = "attached"
		aDrive.DeviceFile = drive.DeviceFile
		aDrive.DriveType = drive.Type
		aDrive.SizeBytes = drive.SizeBytes
		attacheDrivesSN = append(attacheDrivesSN, aDrive.SerialNumber)
		database.DB.Save(&aDrive)

		for _, partition := range drive.Partitions {
			// Check if we know this partition
			log.Tracef("Agent returned partition %v", partition)
			var aPartition models.DrivePartition
			uniqueName := fmt.Sprintf("%s-%s", drive.SerialNumber, partition.Name)
			result := database.DB.Where("name = ?", uniqueName).First(&aPartition)

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// First time we see this partition, create it
				log.Debugf("First time we see partition with this name: %v. Creating a record.", uniqueName)
				aPartition = models.DrivePartition{Name: uniqueName}
				aPartition.Drive = aDrive
				if partition.Label != "" {
					aPartition.Caption = partition.Label
				} else {
					aPartition.Caption = partition.Name
				}
				database.DB.Create(&aPartition)
			}

			if systemPartitions[partition.MountPoint] || systemDevices[partition.DeviceFile] {
				aPartition.Type = "system"
			}

			aPartition.DeviceFile = partition.DeviceFile
			aPartition.SizeBytes = partition.SizeBytes
			aPartition.AvailableBytes = partition.AvailableBytes
			aPartition.MountPoint = partition.MountPoint
			aPartition.Label = partition.Label
			aPartition.IsReadOnly = partition.IsReadOnly
			aPartition.Filesystem = partition.Filesystem
			if partition.MountPoint != "" {
				aPartition.Status = "mounted"
			} else if aPartition.Filesystem != "" && supportedFilesystems[aPartition.Filesystem] {
				aPartition.Status = "unmounted"
			} else {
				aPartition.Status = "unmountable"
			}
			database.DB.Save(&aPartition)
		}
	}
	// Update status of drives no longer attached
	database.DB.Model(&models.Drive{}).Not(map[string]interface{}{"serial_number": attacheDrivesSN}).Update("status", "detached")
	return nil
}

func updateInternalStorageDummyPartition(internalStoragePath string) {
	// Create a dummy partition for "Internal Storage"
	// Internal Storage is a dummy partition that belongs a real mounted partition on some drive
	// Let's find what real drive it belongs to
	var partitions []models.DrivePartition
	database.DB.Joins("Drive").Where("drive.status = ?", "attached").Find(&partitions)

	var internalStorageDrive *models.Drive
	var internalPartitionParent *models.DrivePartition
	var longestPrefix string
	var tries []string
	for _, partition := range partitions {
		tries = append(tries, partition.MountPoint)
		if strings.HasPrefix(internalStoragePath, partition.MountPoint) &&
			len(partition.MountPoint) > len(longestPrefix) {
			// Internal Storage directory lives in this partition, let's create another
			internalStorageDrive = &partition.Drive
			internalPartitionParent = &partition
			longestPrefix = partition.MountPoint
		}
	}
	if internalStorageDrive == nil {
		log.Errorf("Could not determine drive for internal storage path: %s, tried: %v",
			internalStoragePath, tries)
		log.Warningf("Stats for internal storage will be incorrect due to failure in determining its drive")
		internalStorageDrive = &models.Drive{Status: "attached"}
		internalPartitionParent = &models.DrivePartition{}
	} else {
		log.Tracef("Internal storage path: %s belongs to drive: %s and partition: %s",
			internalStoragePath, internalStorageDrive.DeviceFile, internalPartitionParent.Name)
	}

	var internalStoragePartition models.DrivePartition
	database.DB.FirstOrInit(&internalStoragePartition, models.DrivePartition{
		Name:       "internal-storage",
		DeviceFile: "/dev/dummy1",
	})
	internalStoragePartition.AvailableBytes = internalPartitionParent.AvailableBytes
	internalStoragePartition.SizeBytes = internalPartitionParent.SizeBytes
	internalStoragePartition.Drive = *internalStorageDrive
	internalStoragePartition.Filesystem = internalPartitionParent.Filesystem
	internalStoragePartition.MountPoint = internalStoragePath
	internalStoragePartition.Label = internalPartitionParent.Label
	internalStoragePartition.Status = "mounted"
	internalStoragePartition.Caption = "Internal Storage"
	internalStoragePartition.Type = "internal"
	internalStoragePartition.OrderNumber = 0
	database.DB.Save(&internalStoragePartition)
}

func updateIndexingProgress() {
	var partitions []models.DrivePartition
	database.DB.Where("index_status <> ?", "").Find(&partitions)

	// Update index status (if partition has been indexed)
	for _, aPartition := range partitions {
		var indexProgressRes *zahif_pb.IndexProgressResponse
		indexID := fmt.Sprintf("%s-%d", aPartition.Type, aPartition.ID)

		indexProgressRes, err := zahif.Client.IndexProgress(context.Background(), &zahif_pb.IndexProgressRequest{
			IndexIdentifier: indexID,
		})

		if err == nil {
			aPartition.IndexProgress = indexProgressRes.PercentageDone
			if aPartition.IndexProgress >= 100 && (aPartition.IndexStatus == "creating" || aPartition.IndexStatus == "resuming") {
				aPartition.IndexStatus = "created"
			} else if !indexProgressRes.IsRunning && aPartition.IndexStatus == "deleting" {
				aPartition.IndexStatus = "" // index has been fully deleted
			} else if !indexProgressRes.IsRunning && aPartition.IndexStatus == "pausing" {
				aPartition.IndexStatus = "paused"
			} else if indexProgressRes.IsRunning && aPartition.IndexStatus == "resuming" {
				aPartition.IndexStatus = "creating"
			}

			aPartition.IndexTotalDocuments = indexProgressRes.TotalDocuments
			aPartition.IndexIndexedDocuments = indexProgressRes.IndexedDocuments
		} else if strings.Contains(err.Error(), "Index Does Not Exist") && aPartition.IndexStatus == "deleting" {
			aPartition.IndexStatus = ""
			aPartition.IndexProgress = 0
			aPartition.IndexTotalDocuments = 0
			aPartition.IndexIndexedDocuments = 0
		} else {
			log.Errorf("Unexpected error: %v", err)
		}
		database.DB.Save(&aPartition)
	}
}
