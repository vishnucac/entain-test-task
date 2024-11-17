package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"entain-test-task/db"
	"entain-test-task/models"
	"entain-test-task/utils"
	"github.com/gorilla/mux"
)

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userIdParam := mux.Vars(r)["userId"]
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil || userId <= 0 {
		utils.Error("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, userId).Error; err != nil {
		utils.Error(fmt.Sprintf("User %d not found", userId))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	utils.Info(fmt.Sprintf("Balance fetched for user %d", userId))
	response := map[string]interface{}{
		"userId":  user.UserID,
		"balance": user.Balance,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
