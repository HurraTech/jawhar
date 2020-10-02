package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"hurracloud.io/jawhar/internal/agent"
	pb "hurracloud.io/jawhar/internal/agent/proto"
	"hurracloud.io/jawhar/internal/database"
	"hurracloud.io/jawhar/internal/models"
	"hurracloud.io/jawhar/internal/system"
	zahif "hurracloud.io/jawhar/internal/zahif"
	zahif_pb "hurracloud.io/jawhar/internal/zahif/proto"
)

type Controller struct {
	MountPointsRoot      string
	ContainersRoot       string
	InternalStoragePath  string
	SupportedFilesystems map[string]bool
	SouqAPI              string
	SouqUsername         string
	SouqPassword         string
}

/* GET /sources */
func (c *Controller) GetSources(ctx echo.Context) error {

	system.UpdateSources(c.InternalStoragePath)

	var partitions []models.DrivePartition
	database.DB.
		Order("order_number asc").Order("drive_partitions.id asc").Joins("Drive").
		Where("drive.status = ? AND type <> ?", "attached", "system").
		Or(models.DrivePartition{Type: "internal"}).
		Find(&partitions)

	return ctx.JSON(http.StatusOK, partitions)
}

/* POST /sources/:type/:id/mount */
func (c *Controller) MountSource(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	if sourceType == "partition" || sourceType == "internal" {
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

		// let's ask zahif to resume watching for file change events and index them
		if partition.IndexStatus != "" && partition.IndexStatus != "paused" {
			indexID := fmt.Sprintf("%s-%s", sourceType, sourceId)
			log.Debugf("Making Index Request to Zahif: IndexID=%s", indexID)
			_, err := zahif.Client.StartOrResumeIndex(context.Background(), &zahif_pb.IndexRequest{
				IndexIdentifier: indexID,
			})
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("unexpected error: %s", err)})
			}
		}

	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "partition moutned"})
}

/* POST /sources/:type/:id/search */
func (c *Controller) SearchSource(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")

	if ctx.QueryParam("q") == "" || ctx.QueryParam("from") == "" || ctx.QueryParam("to") == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "query string parameters 'from', 'to' and 'q' are required"})
	}

	q := ctx.QueryParam("q")
	from, err := strconv.Atoi(ctx.QueryParam("from"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("Could not parse 'from' Query Parameter: %s: %v", ctx.QueryParam("from"), err)})
	}

	to, err := strconv.Atoi(ctx.QueryParam("to"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("Could not parse 'to' Query Parameter: %s: %v", ctx.QueryParam("to"), err)})
	}

	limit := (to - from + 1)
	results := []map[string]interface{}{}
	indexID := fmt.Sprintf("%s-%s", sourceType, sourceId)

	if sourceType == "partition" || sourceType == "internal" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		log.Debugf("Making Search Request to Zahif: Query=%s, IndexID=%s", q, indexID)
		res, err := zahif.Client.SearchIndex(context.Background(), &zahif_pb.SearchIndexRequest{
			IndexIdentifier: indexID,
			Query:           q,
			Limit:           int32(limit),
			Offset:          int32(from),
		})

		if err != nil {
			log.Errorf("Error while searching zahif: %s", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		log.Debugf("Zahif returned following doucments (%d): %v", len(res.Documents), res.Documents)

		for _, filePath := range res.Documents {
			log.Tracef("Processing result %v", filePath)
			f, err := os.Stat(filePath)
			if os.IsNotExist(err) {
				// file no longer exists
				continue
			}

			file := map[string]interface{}{
				"Name":         f.Name(),
				"Path":         fmt.Sprintf("/sources/%s/%d%s", partition.Type, partition.ID, strings.Replace(filePath, partition.MountPoint, "", 1)),
				"LastModified": f.ModTime(),
				"IsDir":        f.IsDir(),
				"SizeBytes":    f.Size(),
				"Extension":    strings.TrimLeft(path.Ext(f.Name()), "."),
			}
			results = append(results, file)

		}
	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}

	log.Tracef("Reutnring search results: %v", results)

	return ctx.JSON(http.StatusOK, results)
}

/* POST /sources/:type/:id/index */
func (c *Controller) IndexSource(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	indexID := fmt.Sprintf("%s-%s", sourceType, sourceId)
	var excludePatterns []string
	m := echo.Map{}
	if err := ctx.Bind(&m); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "failed to decode payload"})
	}
	log.Debugf("Received IndexSource request with payload: %v", m)

	if patterns, ok := m["excludes"]; ok {
		for _, pattern := range patterns.([]interface{}) {
			excludePatterns = append(excludePatterns, pattern.(string))
		}
		log.Debugf("ExcludePatterns is: %v", excludePatterns)
	}

	if sourceType == "partition" || sourceType == "internal" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		} else if partition.Status != "mounted" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "cannot index unmounted source"})
		}

		log.Debugf("Making Batch Index Request to Zahif: IndexID=%s", indexID)
		_, err := zahif.Client.StartOrResumeIndex(context.Background(), &zahif_pb.IndexRequest{
			IndexIdentifier: indexID,
			Target:          partition.MountPoint,
			ExcludePatterns: excludePatterns,
		})

		if err != nil {
			log.Errorf("Error while creating index on zahif: %s", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		partition.IndexStatus = "creating"
		partition.IndexExcludePatterns = strings.Join(excludePatterns, "|||")
		database.DB.Save(&partition)

	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "index scheduled"})
}

/* DELETE /sources/:type/:id/index */
func (c *Controller) DeleteIndex(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	indexID := fmt.Sprintf("%s-%s", sourceType, sourceId)

	if sourceType == "partition" || sourceType == "internal" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		} else if partition.IndexStatus == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "source is not indexed"})
		}

		log.Debugf("Making Delete Index Request to Zahif: IndexID=%s", indexID)
		_, err := zahif.Client.DeleteIndex(context.Background(), &zahif_pb.DeleteIndexRequest{
			IndexIdentifier: indexID,
		})

		if err != nil {
			log.Errorf("Error while deleting index on zahif: %s", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		partition.IndexStatus = "deleting"
		database.DB.Save(&partition)

	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "index deletion scheduled"})
}

/* POST /sources/:type/:id/resumeIndex */
func (c *Controller) ResumeIndex(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	indexID := fmt.Sprintf("%s-%s", sourceType, sourceId)

	if sourceType == "partition" || sourceType == "internal" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		} else if partition.IndexStatus == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "source is not indexed"})
		}

		log.Debugf("Making Batch Index Request to Zahif: IndexID=%s", indexID)
		_, err := zahif.Client.StartOrResumeIndex(context.Background(), &zahif_pb.IndexRequest{
			IndexIdentifier: indexID,
			Target:          partition.MountPoint,
			ExcludePatterns: strings.Split(partition.IndexExcludePatterns, "|||"),
		})

		if err != nil {
			log.Errorf("Error while resuming index on zahif: %s", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		partition.IndexStatus = "resuming"
		database.DB.Save(&partition)

	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "index resume scheduled"})
}

/* POST /sources/:type/:id/pauseIndex */
func (c *Controller) PauseIndex(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	indexID := fmt.Sprintf("%s-%s", sourceType, sourceId)

	if sourceType == "partition" || sourceType == "internal" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		} else if partition.IndexStatus == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "source is not indexed"})
		}

		log.Debugf("Making Stop Index Request to Zahif: IndexID=%s", indexID)
		_, err := zahif.Client.StopIndex(context.Background(), &zahif_pb.StopIndexRequest{
			IndexIdentifier: indexID,
		})

		if err != nil {
			log.Errorf("Error while stopping index on zahif: %s", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		partition.IndexStatus = "pausing"
		database.DB.Save(&partition)

	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("unsupported type '%s'", sourceType)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "index stop scheduled"})
}

/* POST /sources/:type/:id/unmount */
func (c *Controller) UnmountSource(ctx echo.Context) error {
	sourceType := ctx.Param("type")
	sourceId := ctx.Param("id")
	if sourceType == "partition" || sourceType == "internal" {
		var partition models.DrivePartition
		result := database.DB.Where("id = ?", sourceId).First(&partition)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, result.Error)
		} else if result.Error != nil {
			log.Error("Unexpected error querying DB:", result.Error)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		if partition.IndexStatus == "creating" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "please pause or cancel indexing"})
		} else if partition.IndexStatus != "" {
			// let's ask zahif stop watching for and indexing file changes
			indexID := fmt.Sprintf("%s-%s", sourceType, sourceId)
			log.Debugf("Making Stop Index Request to Zahif: IndexID=%s", indexID)
			_, err := zahif.Client.StopIndex(context.Background(), &zahif_pb.StopIndexRequest{
				IndexIdentifier: indexID,
			})

			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("unexpected error: %s", err)})
			}
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
	if sourceType == "partition" || sourceType == "internal" {
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
	} else if err != nil {
		log.Errorf("Could not stat directory: %s: %s", targetPath, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
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
			log.Errorf("Error while reading directory %s: %s", targetPath, err)
			if os.IsPermission(err) {
				return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Access Denied"})
			} else {
				return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
			}
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

/* GET /app/store */
func (c *Controller) GetStoreApps(ctx echo.Context) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.SouqAPI, "apps"), nil)
	if err != nil {
		log.Errorf("Error connecting Souq API: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	req.SetBasicAuth(c.SouqUsername, c.SouqPassword)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error connecting Souq API: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return ctx.String(http.StatusOK, string(body))
}

/* GET /commands/:id */
func (c *Controller) GetCommand(ctx echo.Context) error {
	var cmd models.AppCommand
	result := database.DB.Where("id = ?", ctx.Param("id")).First(&cmd)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": "command not found"})
	}

	return ctx.JSON(http.StatusOK, cmd)
}

/* GTE /apps */
func (c *Controller) ListInstalledApps(ctx echo.Context) error {
	var apps []models.App
	database.DB.Find(&apps)

	return ctx.JSON(http.StatusOK, apps)
}

/* GET /apps/:id/webapp/* */
func (c *Controller) ProxyWebApp(ctx echo.Context) error {
	var app models.App
	result := database.DB.Preload("WebApp").Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": "app is not installed"})
	}

	targetPort := app.UIPort
	upstreamURL := fmt.Sprintf("http://localhost:%d/%s", targetPort, ctx.Param("*"))
	log.Debugf("Proxying request %s to %s", ctx.Param("*"), upstreamURL)

	client := &http.Client{}
	req := ctx.Request().Clone(ctx.Request().Context())
	req.URL.Path = strings.Replace(req.URL.Path, fmt.Sprintf("/apps/%s/webapp", app.UniqueID), "", 1)
	req.URL.Host = fmt.Sprintf("localhost:%d", targetPort)
	req.URL.Scheme = "http"
	req.RequestURI = ""
	req.Header.Add("X-Forwarded-Host", ctx.Request().Host)
	req.Header.Add("X-Forwarded-For", ctx.Request().RemoteAddr)
	req.Header.Add("X-Real-IP", ctx.Request().RemoteAddr)
	req.Header.Add("X-Forwarded-Proto", ctx.Request().URL.Scheme)
	req.Header.Del("Accept-Encoding")
	req.Header.Add("X-Origin-Host", ctx.Request().Host)

	res, err := client.Do(req)

	if err != nil {
		log.Errorf("Error proxying request: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}
	defer res.Body.Close()

	log.Debugf("Upstream request: %v", req)
	return ctx.Stream(http.StatusOK, res.Header.Get("Content-Type"), res.Body)
}

/* GET /apps/:id */
func (c *Controller) GetApp(ctx echo.Context) error {
	var app models.App
	result := database.DB.Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": "app is not installed"})
	}

	return ctx.JSON(http.StatusOK, app)
}

/* GET /apps/:id/state */
func (c *Controller) GetAppState(ctx echo.Context) error {
	var app models.App
	log.Debugf("Fetching state for app %s", ctx.Param("id"))
	result := database.DB.Preload("State").Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app is not installed"})
	}

	return ctx.String(http.StatusOK, app.State.State)
}

/* POST /apps/:id/state */
func (c *Controller) StoreAppState(ctx echo.Context) error {
	var app models.App
	var state []byte
	result := database.DB.Preload("State").Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app is not installed"})
	}

	defer ctx.Request().Body.Close()
	state, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Errorf("Error reding request body: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(state), &jsonMap)
	if err != nil {
		log.Errorf("Error parsing state: %s", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid JSON: %s", err)})
	}

	log.Debugf("Storing state %s for app %s", state, app.UniqueID)

	app.State.State = string(state)
	database.DB.Save(&app.State)

	return ctx.String(http.StatusOK, app.State.State)
}

/* PATCH /apps/:id/state */
func (c *Controller) PatchAppState(ctx echo.Context) error {
	var app models.App
	var state []byte
	result := database.DB.Preload("State").Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app is not installed"})
	}

	defer ctx.Request().Body.Close()
	state, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Errorf("Error reding request body: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	patchMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(state), &patchMap)
	if err != nil {
		log.Errorf("Error parsing state: %s", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid JSON: %s", err)})
	}

	currentState := make(map[string]interface{})
	err = json.Unmarshal([]byte(app.State.State), &currentState)
	if err != nil {
		log.Errorf("Error parsing current state: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	for k, val := range patchMap {
		currentState[k] = val
	}

	stateStr, err := json.Marshal(currentState)
	if err != nil {
		log.Errorf("Error decoding patched state: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	app.State.State = string(stateStr)
	database.DB.Save(&app.State)

	return ctx.String(http.StatusOK, app.State.State)
}

/* POST /apps/:id/:container/command */
func (c *Controller) ExecAppCommand(ctx echo.Context) error {
	var app models.App
	result := database.DB.Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app is not installed"})
	}

	m := echo.Map{}
	if err := ctx.Bind(&m); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("failed to decode payload: %s", err)})
	}

	if _, ok := m["Env"]; !ok {
		m["Env"] = map[string]interface{}{}
	}
	if _, ok := m["Args"]; !ok {
		m["Args"] = []string{}
	}

	if _, ok := m["Cmd"]; !ok {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required 'Cmd' body parameter"})
	}

	envStr, err := json.Marshal(m["Env"])
	if err != nil {
		log.Errorf("Error parsing command env: %s", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid env JSON: %s", err)})
	}

	argStr, err := json.Marshal(m["Args"])
	if err != nil {
		log.Errorf("Error parsing command args: %s", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("invalid arg JSON: %s", err)})
	}

	appCommand := models.AppCommand{
		App:    app,
		Cmd:    m["Cmd"].(string),
		Args:   string(argStr),
		Env:    string(envStr),
		Status: "running",
	}
	database.DB.Create(&appCommand)

	res, err := agent.Client.ExecInContainerSpec(context.Background(), &pb.ExecInContainerSpecRequest{
		Name:          app.UniqueID,
		Context:       c.ContainersRoot,
		Spec:          app.ContainerSpec,
		ContainerName: ctx.Param("container"),
		Cmd:           appCommand.Cmd,
		Args:          appCommand.Args,
		Env:           appCommand.Env,
	})

	if err != nil {
		log.Errorf("Error executing command '%s' in %s/%s: %s", appCommand.Cmd, app.UniqueID, ctx.Param("container"), err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	appCommand.Status = "completed"
	appCommand.Output = res.Output
	database.DB.Save(&appCommand)

	return ctx.JSON(http.StatusOK, appCommand)
}

/* PUT /apps/:id/:container */
func (c *Controller) StartAppContainer(ctx echo.Context) error {
	var app models.App
	result := database.DB.Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app is not installed"})
	}

	log.Debugf("Starting container: %s/%s", app.UniqueID, ctx.Param("container"))

	_, err := agent.Client.StartContainerInSpec(context.Background(), &pb.ContainerSpecRequest{
		Name:          app.UniqueID,
		Context:       c.ContainersRoot,
		Spec:          app.ContainerSpec,
		ContainerName: ctx.Param("container"),
	})

	if err != nil {
		log.Errorf("Error starting container in %s/%s: %s", app.UniqueID, ctx.Param("container"), err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.JSON(http.StatusOK, app)
}

/* DELETE /apps/:id/:container */
func (c *Controller) StopAppContainer(ctx echo.Context) error {
	var app models.App
	result := database.DB.Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app is not installed"})
	}

	log.Debugf("Stopping container: %s/%s", app.UniqueID, ctx.Param("container"))

	_, err := agent.Client.StopContainerInSpec(context.Background(), &pb.ContainerSpecRequest{
		Name:          app.UniqueID,
		Context:       c.ContainersRoot,
		Spec:          app.ContainerSpec,
		ContainerName: ctx.Param("container"),
	})

	if err != nil {
		log.Errorf("Error stopping container in %s/%s: %s", app.UniqueID, ctx.Param("container"), err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.JSON(http.StatusOK, app)
}

/* POST /apps/:id */
func (c *Controller) InstallApp(ctx echo.Context) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/apps/%s", c.SouqAPI, ctx.Param("id")), nil)
	if err != nil {
		log.Errorf("Error connecting Souq API: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	req.SetBasicAuth(c.SouqUsername, c.SouqPassword)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error connecting Souq API: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error connecting Souq API: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	souqApp := models.App{}
	err = json.Unmarshal(body, &souqApp)
	if err != nil {
		log.Errorf("Error parsing Souq API: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	// Do we have this app already installed?
	var app models.App
	result := database.DB.Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		app = souqApp
		emptyState := models.AppState{State: "{}"}
		app.Status = "installing"
		app.State = emptyState
		database.DB.Create(&emptyState)
		database.DB.Create(&app)
	} else if app.Version == souqApp.Version {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app already installed"})
	} else {
		app.Status = "updating"
		app.Containers = souqApp.Containers
		database.DB.Save(&app)
	}

	go func() {
		// Let's ask agent to load the lateast image
		appImageURL := fmt.Sprintf("%s/apps/%s/image", c.SouqAPI, app.UniqueID)
		log.Debugf("Downloading UI image %s", appImageURL)
		_, err = agent.Client.LoadImage(context.Background(),
			&pb.LoadImageRequest{URL: appImageURL, Username: c.SouqUsername, Password: c.SouqPassword})
		if err != nil {
			log.Errorf("Error loading UI image: %s: %s", appImageURL, err)
			app.Status = "error"
			database.DB.Save(&app)
			return
		}

		// Let's ask agent download and install all dependant container images (if it has any)
		if strings.TrimSpace(app.Containers) != "" {
			for _, image := range strings.Split(app.Containers, ",") {
				log.Tracef("Downloading container image: %s", image)
				_, err = agent.Client.LoadImage(context.Background(),
					&pb.LoadImageRequest{URL: fmt.Sprintf("%s/containers/%s", c.SouqAPI, image), Username: c.SouqUsername, Password: c.SouqPassword})
				if err != nil {
					log.Errorf("Error loading image: %s: %s", image, err)
					app.Status = "error"
					database.DB.Save(&app)
					return
				}
			}

			// Let's retrieve container.yml file
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/apps/%s/containers", c.SouqAPI, app.UniqueID), nil)
			if err != nil {
				log.Errorf("Error retrieving %s/containers.yml: %s", app.UniqueID, err)
				app.Status = "error"
				database.DB.Save(&app)
				return
			}

			req.SetBasicAuth(c.SouqUsername, c.SouqPassword)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Errorf("Error connecting Souq API: %s", err)
				app.Status = "error"
				database.DB.Save(&app)
				return
			}

			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("Error connecting Souq API: %s", err)
				app.Status = "error"
				database.DB.Save(&app)
				return
			}

			log.Tracef("Containers spec for app %s is %s", app.UniqueID, body)
			app.ContainerSpec = string(body)
			database.DB.Save(&app)

			// start sub-containers
			_, err = agent.Client.RunContainerSpec(context.Background(),
				&pb.ContainerSpecRequest{Name: app.UniqueID, Context: c.ContainersRoot, Spec: string(body)})
			if err != nil {
				log.Errorf("Error starting container spec: %s", err)
				app.Status = "error"
				database.DB.Save(&app)
				return
			}
		}

		// Start App UI
		// Find available port
		if app.WebApp.Type == "sdk" {
			listener, err := net.Listen("tcp", ":0")
			if err != nil {
				log.Errorf("Failed to find a free port number: %s", err)
				app.Status = "error"
				database.DB.Save(&app)
				return
			}
			listener.Close() // we're not really using the listener
			app.UIPort = listener.Addr().(*net.TCPAddr).Port
			log.Debugf("Using port %d for UI of app %s", app.UIPort, app.UniqueID)

			_, err = agent.Client.RunContainer(context.Background(),
				&pb.RunContainerRequest{Name: app.UniqueID,
					Image:             fmt.Sprintf("%s:%s", app.UniqueID, app.Version),
					PortMappingSource: uint32(app.UIPort),
					PortMappingTarget: 3000,
					Env:               fmt.Sprintf("REACT_APP_AUID=%s", app.UniqueID),
				})

			if err != nil {
				log.Errorf("Error starting UI container: %s", err)
				app.Status = "error"
				database.DB.Save(&app)
				return
			}
		} else if app.WebApp.Type == "container" {

			res, err := agent.Client.GetContainerPortBindingInSpec(context.Background(),
				&pb.ContainerPortBindingInSpecRequest{
					Name:          app.UniqueID,
					Context:       c.ContainersRoot,
					Spec:          app.ContainerSpec,
					ContainerName: app.WebApp.TargetContainer,
					ContainerPort: uint32(app.WebApp.TargetPort),
				})

			if err != nil {
				log.Errorf("Error determining web app port: %s", err)
				app.Status = "error"
				database.DB.Save(&app)
				return
			}

			app.UIPort = int(res.PortBinding)

		}

		app.Status = "installed"
		database.DB.Save(&app)
	}()

	return ctx.JSON(http.StatusOK, app)
}

/* DELETE /apps/:id */
func (c *Controller) DeleteApp(ctx echo.Context) error {
	var app models.App
	result := database.DB.Where("unique_id = ?", ctx.Param("id")).First(&app)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "app is not installed"})
	} else {
		app.Status = "deleting"
		database.DB.Save(&app)
	}

	go func() {
		// Stop UI container
		_, err := agent.Client.KillContainer(context.Background(),
			&pb.KillContainerRequest{Name: app.UniqueID})

		if err != nil {
			log.Errorf("Error stopping UI container: %s", err)
			app.Status = "error"
			database.DB.Save(&app)
		}

		// Let's ask agent kill and unload all dependant container images (if it has any)
		if strings.TrimSpace(app.Containers) != "" {
			// kill containers
			_, err = agent.Client.RemoveContainerSpec(context.Background(),
				&pb.ContainerSpecRequest{Name: app.UniqueID, Context: c.ContainersRoot, Spec: app.ContainerSpec})
			if err != nil {
				log.Errorf("Error removing container spec for app %s: %s", app.UniqueID, err)
				app.Status = "error"
				database.DB.Save(&app)
			}

			for _, image := range strings.Split(app.Containers, ",") {
				log.Tracef("Unloading container image: %s", image)
				_, err = agent.Client.UnloadImage(context.Background(),
					&pb.UnloadImageRequest{Tag: image})
				if err != nil {
					log.Errorf("Error unloading image: %s: %s", image, err)
					app.Status = "error"
					database.DB.Save(&app)
				}
			}
		}

		// Let's ask remove images to clean up space
		appImageName := fmt.Sprintf("%s:%s", app.UniqueID, app.Version)
		_, err = agent.Client.UnloadImage(context.Background(), &pb.UnloadImageRequest{Tag: appImageName})
		if err != nil {
			log.Errorf("Error unloading image: %s: %s", appImageName, err)
			app.Status = "error"
			database.DB.Save(&app)
		}

		database.DB.Delete(&app)
	}()

	return ctx.JSON(http.StatusOK, app)
}
