package handlers

import (
	"encoding/json"
	"entain-test-task/db"
	"entain-test-task/models"
	"entain-test-task/utils"
	"net/http"
)

func UserListHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		utils.Error("Failed to retrieve users from the database")
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	utils.Info("User list fetched successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
