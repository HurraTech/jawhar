package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hurracloud.io/jawhar/internal/models"
)

var (
	DB *gorm.DB
)

func OpenDatabase(dbFile string, debug bool) {
	var err error
	logLevel := logger.Warn
	if debug {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		panic("failed to open database")
	}
}

func Migrate() {
	DB.AutoMigrate(models.Drive{}, models.DrivePartition{}, models.App{}, models.AppState{}, models.AppCommand{}, models.WebApp{})
}
