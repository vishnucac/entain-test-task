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

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userId from URL
	userIdParam := mux.Vars(r)["userId"]
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil || userId <= 0 {
		utils.Error("Invalid user ID provided", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Decode the transaction from the request
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		utils.Error("Invalid request payload", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the amount
	if !utils.IsValidAmount(transaction.Amount) {
		utils.Error("Invalid amount format", nil)
		http.Error(w, "Invalid amount format - Amount must be a positive number with up to 2 decimal places", http.StatusBadRequest)
		return
	}

	// Check for duplicate transaction ID
	var existingTransaction models.Transaction
	if err := db.DB.Where("transaction_id = ?", transaction.TransactionID).First(&existingTransaction).Error; err == nil {
		utils.Error(fmt.Sprintf("Duplicate transaction ID: %s", transaction.TransactionID), nil)
		http.Error(w, "Duplicate transaction ID", http.StatusConflict)
		return
	}

	// Retrieve the user from the database
	var user models.User
	if err := db.DB.First(&user, userId).Error; err != nil {
		utils.Error(fmt.Sprintf("User %d not found", userId), err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse amount to float64
	parsedAmount := utils.ParseAmount(transaction.Amount)

	// Apply the transaction
	switch transaction.State {
	case "win":
		user.Balance += parsedAmount
	case "lose":
		if user.Balance < parsedAmount {
			utils.Error(fmt.Sprintf("Transaction would result in a negative balance for user %d", userId), nil)
			http.Error(w, "Insufficient balance - account balance cannot be negative", http.StatusBadRequest)
			return
		}
		user.Balance -= parsedAmount
	default:
		utils.Error("Invalid transaction state", nil)
		http.Error(w, "Invalid transaction state", http.StatusBadRequest)
		return
	}

	// Save the user and transaction to the database
	if err := db.DB.Save(&user).Error; err != nil {
		utils.Error(fmt.Sprintf("Failed to update balance for user %d", userId), err)
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Create(&transaction).Error; err != nil {
		utils.Error("Failed to log transaction", err)
		http.Error(w, "Failed to log transaction", http.StatusInternalServerError)
		return
	}

	utils.Info(fmt.Sprintf("Transaction for user %d processed successfully", userId))
	w.WriteHeader(http.StatusOK)
}
