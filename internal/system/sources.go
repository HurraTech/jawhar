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
		attacheDrivesSN = append(attacheDrivesSN, aDrive.SerialNumber)
		aDrive.Status = "attached"
		database.DB.Save(&aDrive)

		for _, partition := range drive.Partitions {
			// Check if we know this partition
			log.Tracef("Agent returned partition %v", partition)
			var aPartition models.DrivePartition
			uniqueName := fmt.Sprintf("%s-%s", drive.SerialNumber, partition.Index)
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
	database.DB.Model(&models.Drive{}).Not(map[string]interface{}{"serial_number": attacheDrivesSN}).Update("status", "detached")

	return partitions, nil
}

func updateInternalStorageDummyPartition(internalStoragePath string, partitions []models.DrivePartition) {
	// Create a dummy partition for "Internal Storage"
	// Internal Storage is a dummy partition that belongs a real mounted partition on some drive
	// Let's find what real drive it belongs to
	var internalStorageDrive models.Drive
	var internalPartitionParent *models.DrivePartition
	var longestPrefix string
	var tries []string
	for _, partition := range partitions {
		tries = append(tries, partition.MountPoint)
		if partition.Status == "mounted" && strings.HasPrefix(options.CmdOptions.MountPointsRoot, options.CmdOptions.InternalStorage) &&
			strings.HasPrefix(internalStoragePath, partition.MountPoint) &&
			len(partition.MountPoint) > len(longestPrefix) {
			// Internal Storage directory lives in this partition
			internalPartitionParent = &partition
			longestPrefix = partition.MountPoint
			log.Debugf("Drive is candidate for internal storage: %v", partition.Drive)
		}
	}
	if internalPartitionParent == nil {
		log.Errorf("Could not determine drive for internal storage path: %s, tried: %v",
			internalStoragePath, tries)
		log.Warningf("Stats for internal storage will be incorrect due to failure in determining its drive")
		database.DB.Where(models.Drive{Status: "attached",
			SerialNumber: "0",
			Name:         "internal",
			DeviceFile:   "/dev/internal",
			OrderNumber:  0,
			Vendor:       "HurraCloud",
		}).FirstOrCreate(&internalStorageDrive)

		internalPartitionParent = &models.DrivePartition{}
	} else {
		database.DB.Where("id = ?", internalPartitionParent.DriveID).First(&internalStorageDrive)
		log.Debugf("Internal storage path: %s belongs to drive: %v and partition: %s",
			internalStoragePath, internalStorageDrive, internalPartitionParent.Name)
	}

	if internalStorageDrive.OrderNumber != 0 || internalStorageDrive.DriveType != "internal" {
		database.DB.Debug().Model(internalStorageDrive).Updates(models.Drive{OrderNumber: -1, DriveType: "internal"})
	}

	var internalStoragePartition models.DrivePartition
	database.DB.Where(models.DrivePartition{
		DeviceFile: "/dev/dummy1",
	}).FirstOrInit(&internalStoragePartition)
	internalStoragePartition.AvailableBytes = internalPartitionParent.AvailableBytes
	internalStoragePartition.SizeBytes = internalPartitionParent.SizeBytes
	internalStoragePartition.DriveID = internalStorageDrive.ID
	internalStoragePartition.Filesystem = internalPartitionParent.Filesystem
	internalStoragePartition.MountPoint = internalStoragePath
	internalStoragePartition.Status = "mounted"
	internalStoragePartition.Caption = "Internal Storage"
	internalStoragePartition.Type = "internal"
	internalStoragePartition.OrderNumber = -1
	database.DB.Debug().Omit("Drive").Save(&internalStoragePartition)
}

func updateIndexingProgress(partitions []models.DrivePartition) {
	database.DB.Where("index_status <> ?", "").Find(&partitions)

	// Update index status (if partition has been indexed)
	for _, aPartition := range partitions {
		var indexProgressRes *zahif_pb.IndexProgressResponse
		indexID := fmt.Sprintf("%s-%d", aPartition.Type, aPartition.Index)

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
