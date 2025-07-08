package handlers

import (
	"Shiso_Checker/db"
	"Shiso_Checker/models"
	"Shiso_Checker/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		Birthday string `json:"birthday"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	bday, err := time.Parse("2006-01-02", input.Birthday)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// パスワードハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		UserID:       input.UserID,
		Birthday:     bday,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         user.ID,
		"user_id":    user.UserID,
		"birthday":   user.Birthday.Format("2006-01-02"),
		"ideology":   user.Ideology,
		"created_at": user.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func UpdateIdeology(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		Ideology models.Ideology `json:"ideology"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Ideology = &payload.Ideology
	db.DB.Save(&user)

	w.WriteHeader(http.StatusNoContent)
}

func GetIdeologyStats(w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(stats)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// ユーザー情報をレスポンス用に整形
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":         user.ID,
		"user_id":    user.UserID,
		"birthday":   user.Birthday.Format("2006-01-02"),
		"ideology":   user.Ideology,
		"created_at": user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 実行前に存在確認（なくても Delete はできるが 404 出したい場合）
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 削除処理
	if err := db.DB.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func DeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if err := db.DB.Exec("DELETE FROM users").Error; err != nil {
		http.Error(w, "Failed to delete users", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
