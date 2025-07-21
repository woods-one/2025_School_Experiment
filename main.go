package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/woods-one/2025_School_Experiment/db"
	"github.com/woods-one/2025_School_Experiment/handlers"
	"github.com/woods-one/2025_School_Experiment/models"
)

func main() {
	// DB 初期化
	db.Init()

	// 自動マイグレーション
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	r := gin.Default()

	// CORS設定（必要に応じて調整）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 認証不要のルート
	r.POST("/login", handlers.Login)
	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetAllUsers)
	r.DELETE("/users", handlers.DeleteAllUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.DELETE("/users/:id", handlers.DeleteUser)
	r.GET("/stats/ideology", handlers.GetIdeologyStats)

	// 認証が必要なルート
	auth := r.Group("/")
	auth.Use(handlers.AuthMiddleware())
	{
		auth.GET("/users/me", handlers.GetCurrentUser)
		auth.PATCH("/users/me/ideology", handlers.UpdateIdeology)
	}

	log.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
