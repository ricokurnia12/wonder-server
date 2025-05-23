package database

import (
	"os"

	"github.com/ricokurnia12/wonder-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectionDb() {
	dbUrl := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(dbUrl))
	if err != nil {
		panic(err)
	}

	// 🚨 Inilah yang penting:
	DB = db

	// Migrasi
	DB.AutoMigrate(&models.Event{}, &models.BlogPost{}, &models.Author{}, &models.Photo{})
}
