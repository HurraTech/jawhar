package options

import "github.com/jessevdk/go-flags"

type Options struct {
	Host            string         `short:"b" long:"host" env:"HOST" description:"Host to bind HTTP server to" default:"127.0.0.1"`
	Port            int            `short:"p" long:"port" env:"PORT" description:"Port to listen HTTP server" default:"5050"`
	Database        flags.Filename `short:"d" long:"db" env:"DB" description:"Database filename" default:"./data/jawhar.db"`
	AgentHost       string         `short:"H" long:"agent_host" env:"AGENT_HOST" description:"Agent Server Host" default:"127.0.0.1"`
	AgentPort       int            `short:"P" long:"agent_port" env:"AGENT_PORT" description:"Agent Server Port" default:"10000"`
	ZahifHost       string         `short:"z" long:"zahif_host" env:"ZAHIF_HOST" description:"Zahif Server Host" default:"127.0.0.1"`
	ZahifPort       int            `short:"o" long:"zahif_port" env:"ZAHIF_PORT" description:"Zahif Server Port" default:"10001"`
	SouqAPI         string         `short:"S" long:"souq_api" env:"SOUQ_API" description:"Souq API Host" default:"https://souq.hurracloud.io"`
	MountPointsRoot string         `short:"m" long:"mount_points_root" env:"MOUNT_POINTS_ROOT" description:"Path under which drives should be mounted" default:"./data/mounts"`
	ContainersRoot  string         `short:"D" long:"containers_root" env:"CONTAINERS_ROOT" description:"Containers root context" default:"./data"`
	InternalStorage string         `short:"s" long:"internal_storage" env:"INTERNAL_STORAGE" description:"Path to use for 'Internal Storage'" default:"./data/storage"`
	Verbose         []bool         `short:"v" long:"verbose" description:"Enable verbose logging"`
}

var CmdOptions Options

func Parse() error {
	_, err := flags.Parse(&CmdOptions)
	return err
}
