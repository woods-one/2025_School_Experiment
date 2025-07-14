package db

import (
	"log"

	"github.com/woods-one/2025_School_Experiment/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自動マイグレーション
	DB.AutoMigrate(&models.User{})
}
