package handlers

import (
	"encoding/json"
	"entain-test-task/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestPingHandler(t *testing.T) {
	utils.InitLogger()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/ping", PingHandler)

	router.ServeHTTP(rr, req)

	// Check if the status code is OK (200)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	// Check if the response body contains the correct JSON
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	expected := map[string]string{"status": "service is up"}
	if response["status"] != expected["status"] {
		t.Errorf("Expected %v but got %v", expected, response)
	}
}
