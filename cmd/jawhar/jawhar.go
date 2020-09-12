package main

import (
    "fmt"
    "github.com/jessevdk/go-flags"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"

    "hurracloud.io/jawhar/internal/agent"
    "hurracloud.io/jawhar/internal/controller"
    "hurracloud.io/jawhar/internal/database"
)

type Options struct {
    Host      string         `short:"h" long:"host" env:"HOST" description:"Host to bind HTTP server to" default:"127.0.0.1"`
    Port      int            `short:"p" long:"port" env:"PORT" description:"Port to listen HTTP server" default:"5050"`
    Database  flags.Filename `short:"d" long:"db" env:"DB" description:"Database filename" default:"jawhar.db"`
    AgentHost string         `short:"H" long:"agent-host" env:"AGENT_HOST" description:"Agent Server Host" default:"127.0.0.1"`
    AgentPort int            `short:"P" long:"agent-port" env:"AGENT_PORT" description:"Agent Server Port" default:"10000"`
}

var options Options

func main() {
    _, err := flags.Parse(&options)

    if err != nil {
        panic(err)
    }

    database.OpenDatabase(string(options.Database))
    database.Migrate()
    agent.Connect(options.AgentHost, options.AgentPort)

    e := echo.New()
    e.GET("/sources", controller.GetSources)
    log.Fatal(e.Start(fmt.Sprintf("%s:%d", options.Host, options.Port)))
}
