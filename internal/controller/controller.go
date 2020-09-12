package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"hurracloud.io/jawhar/internal/agent"
	pb "hurracloud.io/jawhar/internal/agent/proto"
	"hurracloud.io/jawhar/internal/database"
	"hurracloud.io/jawhar/internal/models"
)

type Controller struct {
	MountPointsRoot      string
	SupportedFilesystems map[string]bool
}

/* GET /sources */
func (c *Controller) GetSources(ctx echo.Context) error {

	response, err := agent.Client.GetDrives(context.Background(), &pb.GetDrivesRequest{})
	if err != nil {
		log.Error("Agent Client Failed to call GetDrives: ", err)
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
			aPartition.Filesystem = partition.Filesystem
			if partition.MountPoint != "" {
				aPartition.Status = "mounted"
			} else if aPartition.Filesystem != "" && c.SupportedFilesystems[aPartition.Filesystem] {
				aPartition.Status = "unmounted"
			} else {
				aPartition.Status = "unmountable"
			}

			database.DB.Save(&aPartition)
		}

	}

	var partitions []models.DrivePartition
	database.DB.Preload("Drive").Find(&partitions)
	return ctx.JSON(http.StatusOK, partitions)
}

/* POST /sources/:type/:id/mount */
func (c *Controller) MountSource(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	if sourceType == "partition" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		mountPoint := path.Join(c.MountPointsRoot, partition.Caption)
		log.Debugf("Mounting %s at %s", partition.DeviceFile, mountPoint)
		_, err := agent.Client.MountDrive(context.Background(), &pb.MountDriveRequest{DeviceFile: partition.DeviceFile, MountPoint: mountPoint})
		if err != nil {
			log.Error("Agent Client Failed to call MountDrive: ", err)
			return ctx.JSON(http.StatusServiceUnavailable, map[string]string{"message": "failed to mount"})
		}
	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "partition moutned"})
}

/* POST /sources/:type/:id/unmount */
func (c *Controller) UnmountSource(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	if sourceType == "partition" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		log.Debugf("Unmounting %s", partition.DeviceFile)
		_, err := agent.Client.UnmountDrive(context.Background(), &pb.UnmountDriveRequest{DeviceFile: partition.DeviceFile})
		if err != nil {
			log.Error("Agent Client Failed to call UnmountDrive: ", err)
			return ctx.JSON(http.StatusServiceUnavailable, map[string]string{"message": "failed to mount"})
		}
	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "partition unmoutned"})
}

/* GET /sources/:type/:id/* */
func (c *Controller) BrowseSource(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceID := ctx.Param("id")
	requestedPath := ""
	if len(ctx.ParamValues()) > 2 {
		requestedPath = ctx.ParamValues()[2]
	}
	var targetPath string
	var source models.DrivePartition
	if sourceType == "partition" {
		result := database.DB.Where("id = ?", sourceID).First(&source)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		} else if source.Status != "mounted" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "cannot list unmounted source"})
		}

		targetPath = path.Join(source.MountPoint, requestedPath)
		log.Debugf("List files of source %s at %s", source.Name, targetPath)
	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}

	stat, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": fmt.Sprintf("%s: no such file or direcrory", requestedPath)})
	}

	rel, err := filepath.Rel(source.MountPoint, targetPath)
	if strings.Contains(rel, "..") {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "cannot access files outsdie of drive"})
	}

	switch mode := stat.Mode(); {
	case mode.IsDir():
		if requestedPath != "" && !strings.HasSuffix(requestedPath, "/") {
			// for consistent behavior, always force directory listings to end with trailing slash
			log.Warnf("Request does not contain expected trailing slash")
			return ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/sources/%s/%s/%s/", sourceType, sourceID, requestedPath))
		}

		files, err := ioutil.ReadDir(targetPath)
		if err != nil {
			log.Errorf("Error while reading directory %s: %v", targetPath, err)
		}

		var response []map[string]interface{}

		if requestedPath != "" {
			response = append(response,
				map[string]interface{}{
					"Name":         "..",
					"Path":         fmt.Sprintf("/sources/%s/%d/%s", source.Type, source.ID, strings.TrimRight(filepath.Dir(fmt.Sprintf("%s../", requestedPath)), ".")),
					"IsDir":        true,
					"Extension":    "",
					"SizeBytes":    0,
					"LastModified": "",
				})
		}

		for _, f := range files {
			trailingSlash := "/"
			if !f.IsDir() {
				trailingSlash = ""
			}

			file := map[string]interface{}{
				"Name":         f.Name(),
				"Path":         fmt.Sprintf("/sources/%s/%d/%s%s%s", source.Type, source.ID, requestedPath, f.Name(), trailingSlash),
				"LastModified": f.ModTime(),
				"IsDir":        f.IsDir(),
				"SizeBytes":    f.Size(),
				"Extension":    strings.TrimLeft(path.Ext(f.Name()), "."),
			}
			response = append(response, file)
		}

		return ctx.JSON(http.StatusOK, map[string][]map[string]interface{}{"content": response})

	default:
		return ctx.File(targetPath)
	}
}
