package handlers

import (
	"net/http"
	"os"

	"github.com/woods-one/2025_School_Experiment/db"
	"github.com/woods-one/2025_School_Experiment/models"
	"github.com/woods-one/2025_School_Experiment/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWTの秘密鍵の取得 グローバル変数
var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ログインの入力形式の構造体
type LoginInput struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

// ログイン関数
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "無効な入力"})
		return
	}

	var user models.User
	if err := db.DB.Where("user_id = ?", input.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "ユーザーが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "サーバーエラー"})
		return
	}

	// パスワード比較
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "パスワードが間違っています"})
		return
	}

	// JWT生成
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "トークン生成エラー"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
