package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"entain-test-task/db"
	"entain-test-task/handlers"
	"entain-test-task/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Initialize database with GORM
	utils.Info("Initializing database...")
	db.InitDB()

	// Seed predefined users
	utils.Info("Seeding users...")
	db.SeedUsers()

	r := mux.NewRouter()
	r.HandleFunc("/user/{userId}/transaction", handlers.TransactionHandler).Methods("POST")
	r.HandleFunc("/user/{userId}/balance", handlers.BalanceHandler).Methods("GET")
	r.HandleFunc("/status", handlers.PingHandler).Methods("GET")
	r.HandleFunc("/users", handlers.UserListHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	utils.Info(fmt.Sprintf("Server is running on port %s", port))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
