package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Shiso_Checker/db"
	"Shiso_Checker/models"
	"Shiso_Checker/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// 受信用の構造体：Birthday は string
	var input struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Birthday string `json:"birthday"` // ← ここを string に
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// "2006-01-02" 形式で日付をパース
	bday, err := time.Parse("2006-01-02", input.Birthday)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	user := models.User{
		Email:     input.Email,
		Name:      input.Name,
		Birthday:  bday,
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
