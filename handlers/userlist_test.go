package handlers

import (
	"encoding/json"
	"entain-test-task/db"
	"entain-test-task/models"
	"entain-test-task/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserListHandler(t *testing.T) {
	utils.InitLogger()

	mockDb, mock, _ := sqlmock.New()
	db.DB, _ = gorm.Open(postgres.New(postgres.Config{Conn: mockDb}), &gorm.Config{})

	t.Run("Success", func(t *testing.T) {
		// Mock the query to fetch users
		rows := sqlmock.NewRows([]string{"user_id", "balance"}).
			AddRow(1, 100.50).
			AddRow(2, 200.75)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(rows)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users", UserListHandler)

		router.ServeHTTP(rr, req)

		// Validate the status code
		assert.Equal(t, http.StatusOK, rr.Code)

		// Validate the response body
		var users []models.User
		err := json.NewDecoder(rr.Body).Decode(&users)
		assert.Nil(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, users[0].UserID, uint64(1))
		assert.Equal(t, users[1].UserID, uint64(2))

		// Ensure that all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure to Retrieve Users", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnError(fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users", UserListHandler)

		router.ServeHTTP(rr, req)

		// Validate the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Ensure that all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
