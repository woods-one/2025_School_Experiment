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

// ユーザー作成
// POST /users
func CreateUser(c *gin.Context) {
	// ユーザー作成に必要な情報
	var input struct {
		UserID   string `json:"user_id"`
		Birthday string `json:"birthday"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が不正です"})
		return
	}

	bday, err := time.Parse("2006-01-02", input.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日付の形式が正しくありません"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}

	user := models.User{
		UserID:       input.UserID,
		Birthday:     bday,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースエラー"})
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

// ログインしているユーザーの情報を取得する
// GET /users/me
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

// ユーザーの思想を更新する
// PATCH /users/me/ideology
func UpdateIdeology(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "コンテキスト内のユーザーIDが無効です"})
		return
	}

	// left center rightのどれか
	var payload struct {
		Ideology string `json:"ideology"` // stringで受ける
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力が不正です"})
		return
	}

	if payload.Ideology != "left" && payload.Ideology != "center" && payload.Ideology != "right" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な値です"})
		return
	}

	var user models.User
	if err := db.DB.Where("user_id = ?", userIDStr).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	val := models.Ideology(payload.Ideology)
	user.Ideology = &val

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "値の更新に失敗しました"})
		return
	}

	c.Status(http.StatusNoContent)
}

// 全てのユーザーの思想の統計を取る（年齢別）
// GET /users/stats/ideology
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

// 全ユーザーの情報を取得する
// GET /users
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースエラー"})
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

// IDからユーザーの情報を取得する
// GET /users/:id
func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザーIDが無効です"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
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

// IDからユーザーを削除する
// DELETE /users/:id
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザーIDが無効です"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	if err := db.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーの削除に失敗しました"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ユーザーを全削除
// DELETE /users
func DeleteAllUsers(c *gin.Context) {
	if err := db.DB.Exec("DELETE FROM users").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーの全削除に失敗しました"})
		return
	}
	c.Status(http.StatusNoContent)
}
