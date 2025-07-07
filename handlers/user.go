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
)

// POST /users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	u.CreatedAt = time.Now()

	if err := db.DB.Create(&u).Error; err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// PATCH /users/{id}/ideology
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
		http.Error(w, "Invalid request", http.StatusBadRequest)
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

// GET /stats/ideology
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
