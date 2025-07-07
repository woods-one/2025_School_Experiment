// main.go
package main

import (
	"Shiso_Checker/db"
	"Shiso_Checker/handlers"
	"log"
	"net/http"
)

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	db.Init()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateUser(w, r)
		case http.MethodGet:
			handlers.GetAllUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			handlers.UpdateIdeology(w, r)
		case http.MethodGet:
			handlers.GetUserByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/stats/ideology", handlers.GetIdeologyStats)

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux))
}
