package models

import (
	"gorm.io/gorm"
)

type Drive struct {
	gorm.Model
	OrderNumber  int `gorm:"default:1"`
	Name         string
	SerialNumber string `gorm:"index"`
	Status       string `gorm:"index:dstatus_ptype`
	DeviceFile   string
	DriveType    string
	SizeBytes    uint64
	Vendor       string
	IsRemovable  bool
	Partitions   []DrivePartition
}

type DrivePartition struct {
	gorm.Model
	Index                 uint32
	OrderNumber           int `gorm:"default:1"`
	Name                  string
	Caption               string
	DriveID               uint
	Drive                 Drive
	DeviceFile            string
	Label                 string
	IsReadOnly            bool
	SizeBytes             uint64
	AvailableBytes        uint64
	Filesystem            string
	MountPoint            string
	Status                string
	Type                  string `gorm:"default:partition;index:;index:dstatus_ptype"`
	IndexStatus           string
	IndexProgress         float32
	IndexTotalDocuments   int32
	IndexIndexedDocuments int32
	IndexExcludePatterns  string
}

type App struct {
	gorm.Model
	UniqueID        string `gorm:"index"`
	Name            string
	Description     string
	LongDescription string
	Publisher       string
	Version         string
	Icon            string
	Status          string
	Containers      string
	ContainerSpec   string
	WebApp          WebApp
	UIPort          int
	State           AppState
	Commands        []AppCommand
}

type AppState struct {
	gorm.Model
	AppID uint
	State string
}

type WebApp struct {
	gorm.Model
	AppID           uint
	Type            string
	TargetPort      int
	TargetContainer string
}

type AppCommand struct {
	gorm.Model
	AppID     uint
	App       App
	Cmd       string
	Container string
	Env       string
	Args      string
	Output    string
	Status    string
}
