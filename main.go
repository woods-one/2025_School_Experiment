package main

import (
	"Shiso_Checker/db"
	"Shiso_Checker/handlers"
	"log"
	"net/http"
)

func main() {
	db.Init() // DB接続とマイグレーション

	http.HandleFunc("/users", handlers.CreateUser)
	http.HandleFunc("/users/", handlers.UpdateIdeology)
	http.HandleFunc("/stats/ideology", handlers.GetIdeologyStats)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
