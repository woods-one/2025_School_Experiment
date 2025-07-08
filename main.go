package main

import (
	"log"
	"net/http"

	"Shiso_Checker/db"
	"Shiso_Checker/handlers"
	"Shiso_Checker/models"
)

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	// DB 初期化
	db.Init()

	// テーブル作成（自動マイグレーション）
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// /users（全体操作用：登録、取得、全削除）
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateUser(w, r)
		case http.MethodGet:
			handlers.GetAllUsers(w, r)
		case http.MethodDelete:
			handlers.DeleteAllUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// /users/{id}（個別ユーザー操作：取得、思想更新、削除）
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserByID(w, r)
		case http.MethodPatch:
			handlers.UpdateIdeology(w, r)
		case http.MethodDelete:
			handlers.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 思想統計（年齢別）
	http.HandleFunc("/stats/ideology", handlers.GetIdeologyStats)

	// ログイン
	http.HandleFunc("/login", handlers.Login)

	// サーバー起動
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}
