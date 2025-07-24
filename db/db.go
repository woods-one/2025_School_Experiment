package db

import (
	"log"

	"github.com/woods-one/2025_School_Experiment/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// データベースの初期化処理
func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}

	DB.AutoMigrate(&models.User{})
}
