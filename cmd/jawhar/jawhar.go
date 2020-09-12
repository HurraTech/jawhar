package main

import (
    "fmt"
    "github.com/jessevdk/go-flags"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "path/filepath"

    "hurracloud.io/jawhar/internal/agent"
    "hurracloud.io/jawhar/internal/controller"
    "hurracloud.io/jawhar/internal/database"
)

type Options struct {
    Host            string         `short:"h" long:"host" env:"HOST" description:"Host to bind HTTP server to" default:"127.0.0.1"`
    Port            int            `short:"p" long:"port" env:"PORT" description:"Port to listen HTTP server" default:"5050"`
    Database        flags.Filename `short:"d" long:"db" env:"DB" description:"Database filename" default:"jawhar.db"`
    AgentHost       string         `short:"H" long:"agent_host" env:"AGENT_HOST" description:"Agent Server Host" default:"127.0.0.1"`
    AgentPort       int            `short:"P" long:"agent_port" env:"AGENT_PORT" description:"Agent Server Port" default:"10000"`
    MountPointsRoot string         `short:"m" long:"mount_points_root" env:"MOUNT_POINTS_ROOT" description:"Path under which drives should be mounted" default:"./mounts"`
    Verbose         bool           `short:"v" long:"verbose" description:"Enable verbose logging"`
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

    if options.Verbose {
        log.SetLevel(log.DebugLevel)
    }

    database.OpenDatabase(string(options.Database))
    database.Migrate()
    agent.Connect(options.AgentHost, options.AgentPort)

    mountRoot, err := filepath.Abs(options.MountPointsRoot)
    if err != nil {
        log.Warnf("Could not determine absolute path for mounts directory '%s': %s", options.MountPointsRoot, err)
        mountRoot = options.MountPointsRoot
    }
    controller := &controller.Controller{MountPointsRoot: mountRoot, SupportedFilesystems: supportedFilesystems}
    e := echo.New()
    e.GET("/sources", controller.GetSources)
    e.POST("/sources/:type/:id/mount", controller.MountSource)
    e.POST("/sources/:type/:id/unmount", controller.UnmountSource)
    e.GET("/sources/:type/:id", controller.BrowseSource)
    e.GET("/sources/:type/:id/*", controller.BrowseSource)
    log.Fatal(e.Start(fmt.Sprintf("%s:%d", options.Host, options.Port)))
}
