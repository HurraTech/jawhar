package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"hurracloud.io/jawhar/internal/models"
)

var (
	DB *gorm.DB
)

func OpenDatabase(dbFile string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}
}

func Migrate() {
	DB.AutoMigrate(models.Drive{}, models.DrivePartition{}, models.App{}, models.AppState{}, models.AppCommand{}, models.WebApp{})
}
