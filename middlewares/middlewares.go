package middlewares

import (
	"fmt"
	"net/http"
)

// Valid source types for /user/{userId}/transaction
var validSourceTypes = map[string]bool{
	"game":    true,
	"server":  true,
	"payment": true,
}

// CheckSourceTypeForTransaction middleware function for /user/{userId}/transaction route
func CheckSourceTypeForTransaction(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceType := r.Header.Get("Source-Type")

		// If no Source-Type header, return an error
		if sourceType == "" {
			http.Error(w, "Missing Source-Type header", http.StatusBadRequest)
			return
		}

		// Validate Source-Type
		if !validSourceTypes[sourceType] {
			http.Error(w, fmt.Sprintf("Invalid Source-Type: %s", sourceType), http.StatusBadRequest)
			return
		}

		// Proceed to the next handler if Source-Type is valid
		next(w, r)
	}
}
