package database

import (
	"os"
	"fmt"
	"log"
	"time"

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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		  SlowThreshold:              time.Second,   // Slow SQL threshold
		  LogLevel:                   logLevel, // Log level
		  IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
		  ParameterizedQueries:      false,           // Don't include params in the SQL log
		  Colorful:                  true,          // Disable color
		},
	  )

	DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?cache=shared", dbFile)), &gorm.Config{
		Logger: newLogger,
	})

	DB.Exec("PRAGMA journal_mode=WAL;")

	if err != nil {
		panic("failed to open database")
	}
}

func Migrate() {
	DB.AutoMigrate(models.Drive{},
		models.DrivePartition{},
		models.App{},
		models.AppState{},
		models.AppCommand{},
		models.WebApp{},
		models.SystemState{})
}
