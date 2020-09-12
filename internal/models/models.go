package models

import (
	"gorm.io/gorm"
)

type Source interface {
	Type()
}

type Drive struct {
	gorm.Model
	Name         string
	SerialNumber string
	DeviceFile   string
	DriveType    string
	SizeBytes    uint64
	IsRemovable  bool
	Partitions   []DrivePartition
}

type DrivePartition struct {
	gorm.Model
	Name           string
	Caption        string
	DriveID        int
	Drive          Drive
	DeviceFile     string
	Label          string
	IsReadOnly     bool
	SizeBytes      uint64
	AvailableBytes uint64
	Filesystem     string
	MountPoint     string
	Status         string
	Type           string `gorm:"default:partition"`
	Index          Index
}

type Index struct {
	gorm.Model
	DrivePartitionID int
	SizeBytes        int
}
