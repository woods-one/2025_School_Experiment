package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/woods-one/2025_School_Experiment/db"
	"github.com/woods-one/2025_School_Experiment/models"
	"github.com/woods-one/2025_School_Experiment/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var input struct {
		UserID   string `json:"user_id"`
		Birthday string `json:"birthday"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	bday, err := time.Parse("2006-01-02", input.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		UserID:       input.UserID,
		Birthday:     bday,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	response := map[string]interface{}{
		"id":         user.ID,
		"user_id":    user.UserID,
		"birthday":   user.Birthday.Format("2006-01-02"),
		"ideology":   user.Ideology,
		"created_at": user.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func GetCurrentUser(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーIDが見つかりません"})
		return
	}

	userID := userIDVal.(string)

	var user models.User
	if err := db.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"user_id":    user.UserID,
		"birthday":   user.Birthday.Format("2006-01-02"),
		"ideology":   user.Ideology,
		"created_at": user.CreatedAt,
	})
}

func UpdateIdeology(c *gin.Context) {
	userIDInterface, exists := c.Get("userID") // ミドルウェアで設定されたuserIDを取得
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDInterface.(uint)

	var payload struct {
		Ideology models.Ideology `json:"ideology"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Ideology = &payload.Ideology
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ideology"})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetIdeologyStats(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)

	stats := map[string]map[models.Ideology]int{}

	for _, u := range users {
		if u.Ideology == nil {
			continue
		}
		age := utils.GetAge(u.Birthday)
		group := utils.GetAgeGroup(age)

		if stats[group] == nil {
			stats[group] = map[models.Ideology]int{}
		}
		stats[group][*u.Ideology]++
	}

	c.JSON(http.StatusOK, stats)
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var responses []map[string]interface{}
	for _, u := range users {
		responses = append(responses, map[string]interface{}{
			"id":         u.ID,
			"user_id":    u.UserID,
			"birthday":   u.Birthday.Format("2006-01-02"),
			"ideology":   u.Ideology,
			"created_at": u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := map[string]interface{}{
		"id":         user.ID,
		"user_id":    user.UserID,
		"birthday":   user.Birthday.Format("2006-01-02"),
		"ideology":   user.Ideology,
		"created_at": user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteAllUsers(c *gin.Context) {
	if err := db.DB.Exec("DELETE FROM users").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete users"})
		return
	}
	c.Status(http.StatusNoContent)
}
