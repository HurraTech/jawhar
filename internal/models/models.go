package models

import (
	"gorm.io/gorm"
)

type Drive struct {
	gorm.Model
	Name         string
	SerialNumber string
	Status       string
	DeviceFile   string
	DriveType    string
	SizeBytes    uint64
	IsRemovable  bool
	Partitions   []DrivePartition
}

type DrivePartition struct {
	gorm.Model
	Name                  string
	Caption               string
	DriveID               int
	Drive                 Drive
	DeviceFile            string
	Label                 string
	IsReadOnly            bool
	SizeBytes             uint64
	AvailableBytes        uint64
	Filesystem            string
	MountPoint            string
	Status                string
	Type                  string `gorm:"default:partition"`
	IndexStatus           string
	IndexProgress         float32
	IndexTotalDocuments   int32
	IndexIndexedDocuments int32
	IndexExcludePatterns  string
}

type App struct {
	gorm.Model
	UniqueID        string
	Name            string
	Description     string
	LongDescription string
	Publisher       string
	Version         string
	Icon            string
	Status          string
	Containers      string
	ContainerSpec   string
	UIPort          int
}
