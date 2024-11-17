package handlers

import (
	"encoding/json"
	"entain-test-task/utils"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	utils.Info("Ping request received")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "service is up"})
	utils.Info("Ping response sent")
}
