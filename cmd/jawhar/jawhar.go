package main

import (
	"fmt"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"hurracloud.io/jawhar/internal/agent"
	"hurracloud.io/jawhar/internal/controller"
	"hurracloud.io/jawhar/internal/database"
	"hurracloud.io/jawhar/internal/zahif"
)

type Options struct {
	Host            string         `short:"h" long:"host" env:"HOST" description:"Host to bind HTTP server to" default:"127.0.0.1"`
	Port            int            `short:"p" long:"port" env:"PORT" description:"Port to listen HTTP server" default:"5050"`
	Database        flags.Filename `short:"d" long:"db" env:"DB" description:"Database filename" default:"./data/jawhar.db"`
	AgentHost       string         `short:"H" long:"agent_host" env:"AGENT_HOST" description:"Agent Server Host" default:"127.0.0.1"`
	AgentPort       int            `short:"P" long:"agent_port" env:"AGENT_PORT" description:"Agent Server Port" default:"10000"`
	ZahifHost       string         `short:"z" long:"zahif_host" env:"ZAHIF_HOST" description:"Zahif Server Host" default:"127.0.0.1"`
	ZahifPort       int            `short:"o" long:"zahif_port" env:"ZAHIF_PORT" description:"Zahif Server Port" default:"10001"`
	SouqAPI         string         `short:"s" long:"souq_host" env:"SOUQ_API" description:"Souq API Host" default:"http://127.0.0.1:5060"`
	MountPointsRoot string         `short:"m" long:"mount_points_root" env:"MOUNT_POINTS_ROOT" description:"Path under which drives should be mounted" default:"./data/mounts"`
	ContainersRoot  string         `short:"D" long:"containers_root" env:"containers_root" description:"Containers root context" default:"./data"`
	Verbose         []bool         `short:"v" long:"verbose" description:"Enable verbose logging"`
}

var options Options

var supportedFilesystems = map[string]bool{
	"vfat": true,
	"ext4": true,
	"ext3": true,
	"ntfs": true,
}

func main() {
	_, err := flags.Parse(&options)

	if err != nil {
		panic(err)
	}

	if len(options.Verbose) == 1 {
		log.SetLevel(log.DebugLevel)
	} else if len(options.Verbose) > 1 {
		log.SetLevel(log.TraceLevel)
	}

	database.OpenDatabase(string(options.Database))
	database.Migrate()
	agent.Connect(options.AgentHost, options.AgentPort)
	zahif.Connect(options.ZahifHost, options.ZahifPort)

	mountRoot, err := filepath.Abs(options.MountPointsRoot)
	if err != nil {
		log.Warnf("Could not determine absolute path for mounts directory '%s': %s", options.MountPointsRoot, err)
		mountRoot = options.MountPointsRoot
	}

	containersRoot, err := filepath.Abs(options.ContainersRoot)
	if err != nil {
		log.Warnf("Could not determine absolute path for containers directory '%s': %s", options.ContainersRoot, err)
		containersRoot = options.ContainersRoot
	}

	controller := &controller.Controller{MountPointsRoot: mountRoot,
		ContainersRoot:       containersRoot,
		SupportedFilesystems: supportedFilesystems,
		SouqAPI:              options.SouqAPI}
	e := echo.New()
	e.GET("/sources", controller.GetSources)
	e.POST("/sources/:type/:id/mount", controller.MountSource)
	e.POST("/sources/:type/:id/unmount", controller.UnmountSource)
	e.POST("/sources/:type/:id/search", controller.SearchSource)
	e.POST("/sources/:type/:id/index", controller.IndexSource)
	e.DELETE("/sources/:type/:id/index", controller.DeleteIndex)
	e.POST("/sources/:type/:id/pauseIndex", controller.PauseIndex)
	e.POST("/sources/:type/:id/resumeIndex", controller.ResumeIndex)
	e.GET("/sources/:type/:id", controller.BrowseSource)
	e.GET("/sources/:type/:id/*", controller.BrowseSource)
	e.GET("/apps/store", controller.GetStoreApps)
	e.GET("/apps", controller.ListInstalledApps)
	e.GET("/apps/:id", controller.GetApp)
	e.POST("/apps/:id", controller.InstallApp)
	e.DELETE("/apps/:id", controller.DeleteApp)
	log.Fatal(e.Start(fmt.Sprintf("%s:%d", options.Host, options.Port)))
}
