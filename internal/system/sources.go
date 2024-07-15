package system

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"hurracloud.io/jawhar/cmd/jawhar/options"
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

var systemDevicesPrefixes = map[string]bool{
	// RPI
	"/dev/mmcblk0p": true,

	// Thinkpad
	"/dev/nvme0n1p": true,

	// EC2
	"/dev/xvda": true,

	// Mac
	"/dev/disk0s": true,
	"/dev/disk1s": true,
}

var supportedFilesystems = map[string]bool{
	"vfat": true,
	"ext4": true,
	"ext3": true,
	"ntfs": true,
}

func UpdateSources() error {

	partitions, err := updateSources()
	if err != nil {
		return err
	}

	updateInternalStorageDummyPartition(options.CmdOptions.InternalStorage, partitions)
	updateIndexingProgress(partitions)

	return nil
}

func updateSources() ([]models.DrivePartition, error) {
	response, err := agent.Client.GetDrives(context.Background(), &pb.GetDrivesRequest{})
	if err != nil {
		return nil, fmt.Errorf("Agent Client Failed to call GetDrives: %s", err)
	}

	var attacheDrivesSN []string
	var partitions []models.DrivePartition

	if log.GetLevel() == log.TraceLevel {
		log.Tracef("Drive map from hurra agent:")
		for _, drive := range response.Drives {
			log.Tracef(" - %s (SN: %s, controller: %s, type: %s, is_removable: %v)", drive.DeviceFile, drive.SerialNumber,
				drive.StorageController, drive.Type, drive.IsRemovable)
			for _, partition := range drive.Partitions {
				log.Tracef("    |- %s (type: %s)", partition.DeviceFile, partition.Filesystem)
			}
		}
	}
	
	for _, drive := range response.Drives {
		// Check if we know this drive
		log.Tracef("Agent returned drive %v", drive)
		var aDrive models.Drive
		database.DB.FirstOrInit(&aDrive, &models.Drive{SerialNumber: drive.SerialNumber})
		aDrive.Name = drive.Name
		aDrive.DeviceFile = drive.DeviceFile
		aDrive.DriveType = drive.Type
		aDrive.SizeBytes = drive.SizeBytes
		aDrive.Vendor = drive.Vendor
		aDrive.IsRemovable = drive.IsRemovable
		attacheDrivesSN = append(attacheDrivesSN, aDrive.SerialNumber)
		aDrive.Status = "attached"
		database.DB.Save(&aDrive)

		for _, partition := range drive.Partitions {
			// Check if we know this partition
			log.Tracef("Agent returned partition %v", partition)
			var aPartition models.DrivePartition
			uniqueName := fmt.Sprintf("%s-%d", drive.SerialNumber, partition.Index)
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

			if systemPartitions[partition.MountPoint] {
				aPartition.Type = "system"
			} else {
				for prefix, _ := range systemDevicesPrefixes {
					if strings.HasPrefix(partition.DeviceFile, prefix) {
						aPartition.Type = "system"
					}
				}
			}

			aPartition.DeviceFile = partition.DeviceFile
			aPartition.SizeBytes = partition.SizeBytes
			aPartition.AvailableBytes = partition.AvailableBytes
			aPartition.MountPoint = partition.MountPoint
			aPartition.Label = partition.Label
			aPartition.IsReadOnly = partition.IsReadOnly
			aPartition.Filesystem = partition.Filesystem
			var newStatus string
			if partition.MountPoint != "" {
				newStatus = "mounted"
			} else if aPartition.Filesystem != "" && supportedFilesystems[aPartition.Filesystem] {
				newStatus = "unmounted"
			} else {
				newStatus = "unmountable"
			}

			aPartition.Status = newStatus
			database.DB.Save(&aPartition)

			partitions = append(partitions, aPartition)
		}
	}
	// Update status of drives no longer attached
	database.DB.Model(&models.Drive{}).
		Where(map[string]interface{}{"status": "attached"}).
		Not(map[string]interface{}{"serial_number": attacheDrivesSN}).Update("status", "detached")

	return partitions, nil
}

func updateInternalStorageDummyPartition(internalStoragePath string, partitions []models.DrivePartition) {
	// Create a dummy partition for "Internal Storage"
	// Internal Storage is a dummy partition that belongs a real mounted partition on some drive
	// Let's find what real drive it belongs to
	var osDrive models.Drive
	var osPartition models.DrivePartition
	var longestPrefix string
	var tries []string
	log.Debugf("Attempting to determine OS partition and device")
	for _, partition := range partitions {
		tries = append(tries, partition.MountPoint)
		if partition.Status == "mounted" && strings.HasPrefix(internalStoragePath, partition.MountPoint) &&
			len(partition.MountPoint) > len(longestPrefix) {
			// Internal Storage directory lives in this partition
			osPartition = partition
			longestPrefix = partition.MountPoint
			log.Debugf("%v on %v is an OS partition candidate (drive_id=%d, status=%s, length=%d, longestPrefix=%d, hasCorrectPrefix=%s)", 
						partition.MountPoint, partition.DeviceFile, partition.DriveID, partition.Status, len(partition.MountPoint), len(longestPrefix), strings.HasPrefix(internalStoragePath, partition.MountPoint))
		} else {
			log.Debugf("%v is NOT an OS partition candidate (status=%s, length=%d, longestPrefix=%d, hasCorrectPrefix=%s)", 
						partition.MountPoint, partition.Status, len(partition.MountPoint), len(longestPrefix), strings.HasPrefix(internalStoragePath, partition.MountPoint))
		}
	}
	if &osPartition == nil {
		log.Warnf("Could not determine OS partition")
		return
	} else {
		log.Debugf("OS partition is: %v", osPartition.DeviceFile)
		database.DB.Where("id = ?", osPartition.DriveID).First(&osDrive)
		log.Debugf("OS partition: %s belongs to drive: %v",
			osPartition.DeviceFile, osDrive.DeviceFile)
	}

	if osDrive.OrderNumber != 0 || osDrive.DriveType != "internal" {
		database.DB.Model(osDrive).Updates(models.Drive{OrderNumber: -1, DriveType: "internal"})
	}

	osDrive.IsOS = true
	osPartition.IsOS = true
	database.DB.Save(&osPartition)
	database.DB.Save(&osDrive)
}

func updateIndexingProgress(partitions []models.DrivePartition) {
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
			log.Debugf("Index %s progress is %f", indexID, indexProgressRes.PercentageDone)
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
			log.Errorf("Unexpected error while checking index %s progress: %v", indexID, err)
		}
		database.DB.Save(&aPartition)
	}
}
