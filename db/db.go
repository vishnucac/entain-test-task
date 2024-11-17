package db

import (
	"fmt"
	"log"
	"os"

	"entain-test-task/models"
	"entain-test-task/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize database connection
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Transaction{})
	if err != nil {
		utils.Error(fmt.Sprintf("Failed to auto-migrate database: %v", err))
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	utils.Info("Database connected and migrated successfully!")
}

// Seed predefined users in the database
func SeedUsers() {
	users := []models.User{
		{UserID: 1, Balance: 100.00},
		{UserID: 2, Balance: 50.50},
		{UserID: 3, Balance: 200.75},
	}

	for _, user := range users {
		// Check if user already exists
		if err := DB.FirstOrCreate(&user, models.User{UserID: user.UserID}).Error; err != nil {
			utils.Error(fmt.Sprintf("Failed to seed user %d: %v", user.UserID, err))
		} else {
			utils.Info(fmt.Sprintf("User %d seeded successfully", user.UserID))
		}
	}
}
