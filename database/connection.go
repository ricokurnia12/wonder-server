package database

import (
	"github.com/ricokurnia12/wonder-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectionDb() {
	dbUrl := "postgresql://wonderpalembang_owner:npg_Bm7vctlIQ8Wy@ep-icy-mud-a4qi9s46-pooler.us-east-1.aws.neon.tech/wonderpalembang?sslmode=require"

	db, err := gorm.Open(postgres.Open(dbUrl))
	if err != nil {
		panic(err)
	}

	// 🚨 Inilah yang penting:
	DB = db

	// Migrasi
	DB.AutoMigrate(&models.Event{}, &models.BlogPost{}, &models.Author{}, &models.Photo{})
}
