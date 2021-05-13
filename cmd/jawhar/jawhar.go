package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"

	"hurracloud.io/jawhar/cmd/jawhar/options"
	"hurracloud.io/jawhar/internal/agent"
	"hurracloud.io/jawhar/internal/controller"
	"hurracloud.io/jawhar/internal/database"
	"hurracloud.io/jawhar/internal/system"
	"hurracloud.io/jawhar/internal/zahif"
)

func main() {

	if err := options.Parse(); err != nil {
		panic(err)
	}

	if len(options.CmdOptions.Verbose) == 1 {
		log.SetLevel(log.DebugLevel)
	} else if len(options.CmdOptions.Verbose) > 1 {
		log.SetLevel(log.TraceLevel)
	}

	database.OpenDatabase(string(options.CmdOptions.Database), len(options.CmdOptions.Verbose) > 0)
	database.Migrate()
	agent.Connect(options.CmdOptions.AgentHost, options.CmdOptions.AgentPort)
	zahif.Connect(options.CmdOptions.ZahifHost, options.CmdOptions.ZahifPort)

	mountRoot, err := filepath.Abs(options.CmdOptions.MountPointsRoot)
	if err != nil {
		log.Warnf("Could not determine absolute path for mounts directory '%s': %s", options.CmdOptions.MountPointsRoot, err)
		mountRoot = options.CmdOptions.MountPointsRoot
	}

	containersRoot, err := filepath.Abs(options.CmdOptions.ContainersRoot)
	if err != nil {
		log.Warnf("Could not determine absolute path for containers directory '%s': %s", options.CmdOptions.ContainersRoot, err)
		containersRoot = options.CmdOptions.ContainersRoot
	}

	internalStorageAbs, err := filepath.Abs(options.CmdOptions.InternalStorage)
	if err != nil {
		log.Warnf("Could not determine absolute path for internal storage directory '%s': %s", options.CmdOptions.InternalStorage, err)
		internalStorageAbs = options.CmdOptions.InternalStorage
	}
	options.CmdOptions.InternalStorage = internalStorageAbs

	if _, err := os.Stat(internalStorageAbs); os.IsNotExist(err) {
		err := os.MkdirAll(internalStorageAbs, 0755)
		if err != nil {
			log.Fatalf("Could not create internal storage directory: %s: %s", internalStorageAbs, err)
		}
	}

	err = system.UpdateSources()
	if err != nil {
		log.Errorf("Could not create refresh sources: %s", err)
	}

	// Update whenever blocks storage is changed
	w := watcher.New()
	w.AddRecursive("/sys/block")
	w.FilterOps(watcher.Create, watcher.Remove)
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Infof("USB device event %s", event)
				time.Sleep(time.Second)
				system.UpdateSources()
			case <-ticker.C:
				system.UpdateSources()
			case err := <-w.Error:
				log.Errorf("USB Watcher Error: %s", err)
			case <-w.Closed:
				ticker.Stop()
				return
			}
		}
	}()

	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Errorf("Error starting watcher: %s", err)
		}
	}()

	controller := &controller.Controller{MountPointsRoot: mountRoot,
		ContainersRoot: containersRoot,
		SouqAPI:        options.CmdOptions.SouqAPI,
		SouqUsername:   "HURRANET",
		SouqPassword:   "bSdh~e9J:FTbLS#w",
	}
	e := echo.New()

	// storage data
	e.GET("/sources", controller.GetSources)
	e.GET("/partitions", controller.GetPartitions)
	e.POST("/sources/:type/:id/mount", controller.MountSource)
	e.POST("/sources/:type/:id/unmount", controller.UnmountSource)
	e.POST("/sources/:type/:id/search", controller.SearchSource)
	e.POST("/sources/:type/:id/index", controller.IndexSource)
	e.DELETE("/sources/:type/:id/index", controller.DeleteIndex)
	e.POST("/sources/:type/:id/pauseIndex", controller.PauseIndex)
	e.POST("/sources/:type/:id/resumeIndex", controller.ResumeIndex)
	e.GET("/sources/:type/:id", controller.BrowseSource)
	e.GET("/sources/:type/:id/*", controller.BrowseSource)
	e.POST("/sources/:type/:id/*", controller.UploadToSource)
	e.DELETE("/sources/:type/:id/*", controller.DeleteFromSource)

	// Apps
	e.GET("/apps/store", controller.GetStoreApps)
	e.GET("/apps", controller.ListInstalledApps)
	e.GET("/apps/:id", controller.GetApp)
	e.GET("/apps/:id/state", controller.GetAppState)
	e.POST("/apps/:id/state", controller.StoreAppState)
	e.PATCH("/apps/:id/state", controller.PatchAppState)
	e.POST("/apps/:id", controller.InstallApp)
	e.POST("/apps/:id/:container/command", controller.ExecAppCommand)
	e.GET("/commands/:id", controller.GetCommand)
	e.DELETE("/apps/:id", controller.DeleteApp)
	e.PUT("/apps/:id/:container", controller.StartAppContainer)
	e.DELETE("/apps/:id/:container", controller.StopAppContainer)

	// web apps reverse proxies
	e.GET("/apps/:id/webapp/*", controller.ProxyWebApp)
	e.PUT("/apps/:id/webapp/*", controller.ProxyWebApp)
	e.POST("/apps/:id/webapp/*", controller.ProxyWebApp)
	e.DELETE("/apps/:id/webapp/*", controller.ProxyWebApp)

	// system management
	e.GET("/system/stats", controller.GetSystemStats)
	e.POST("/system/update/:version", controller.UpdateSystem)

	log.Fatal(e.Start(fmt.Sprintf("%s:%d", options.CmdOptions.Host, options.CmdOptions.Port)))
}
