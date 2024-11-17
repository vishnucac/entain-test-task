package utils

import (
	"log"
	"regexp"
	"strconv"
)

// isValidAmount checks if the amount is a valid string with up to 2 decimal places.
func IsValidAmount(amount string) bool {
	re := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	return re.MatchString(amount)
}

// ParseAmount converts a valid amount string to a float64.
func ParseAmount(amount string) float64 {
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Printf("Error parsing amount: %v", err)
		return 0 // Handle parsing error safely
	}
	return val
}
