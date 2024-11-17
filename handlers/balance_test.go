package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"entain-test-task/db"
	"entain-test-task/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBalanceHandler(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	dbMock, _ := gorm.Open(dialector, &gorm.Config{})
	db.DB = dbMock

	utils.InitLogger()

	handler := func(w http.ResponseWriter, r *http.Request) {
		BalanceHandler(w, r)
	}

	t.Run("Invalid User ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/invalid/balance", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/{userId}/balance", handler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := 1

		// Mock the query to return no user found
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."user_id" = $1 ORDER BY "users"."user_id" LIMIT $2`)).
			WithArgs(userID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}))

		// Create a request with valid user ID
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d/balance", userID), nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/{userId}/balance", handler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("User Found", func(t *testing.T) {
		var userID float64 = 1
		var mockBalance float64 = 100

		// Mock the query to return the user with a balance
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."user_id" = $1 ORDER BY "users"."user_id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).AddRow(userID, mockBalance))

		// Create a request with valid user ID
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%v/balance", userID), nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/{userId}/balance", handler)
		router.ServeHTTP(rr, req)

		// Validate response code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
		}

		// Validate response body
		var response map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if response["userId"] != userID || response["balance"] != mockBalance {
			t.Errorf("Expected userId %v and balance %v, but got userId %v and balance %v", userID, mockBalance, response["userId"], response["balance"])
		}
	})

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}
