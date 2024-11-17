package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"entain-test-task/db"
	"entain-test-task/models"
	"entain-test-task/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestTransactionHandler(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	dbMock, _ := gorm.Open(dialector, &gorm.Config{})
	db.DB = dbMock

	utils.InitLogger()

	handler := func(w http.ResponseWriter, r *http.Request) {
		TransactionHandler(w, r)
	}

	t.Run("Invalid User ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user/invalid/transaction", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/{userId}/transaction", handler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Transaction, Valid User ID and Transaction", func(t *testing.T) {
		userID := 1
		transaction := models.Transaction{
			TransactionID: "1",
			Amount:        "10",
			State:         "lose",
		}

		// Prepare the mock queries
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "transactions" WHERE transaction_id = $1 ORDER BY "transactions"."transaction_id" LIMIT $2`)).
			WithArgs(transaction.TransactionID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."user_id" = $1 ORDER BY "users"."user_id" LIMIT $2`)).
			WithArgs(userID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).AddRow(userID, 100))

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "balance"=$1 WHERE "user_id" = $2`)).
			WithArgs(90.0, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "transactions" ("transaction_id","amount","state") VALUES ($1,$2,$3)`)).
			WithArgs("1", transaction.Amount, transaction.State).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		txBody, _ := json.Marshal(transaction)

		url := fmt.Sprintf("/user/%d/transaction", userID)
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(txBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/{userId}/transaction", handler)
		router.ServeHTTP(rr, req)

		// Validate response code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
		}

		// Ensure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %v", err)
		}
	})

}
