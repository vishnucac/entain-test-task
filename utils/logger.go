package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	// Create or open the log file
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Log to both console and file
	Logger = log.New(file, "APP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs informational messages
func Info(msg string) {
	Logger.Println("[INFO]:", msg)
}

// Error logs error messages along with the optional error details
func Error(msg string, err ...error) {
	if len(err) > 0 && err[0] != nil {
		Logger.Printf("[ERROR]: %s: %v", msg, err[0])
	} else {
		Logger.Println("[ERROR]:", msg)
	}
}

// Fatal logs fatal error messages and terminates the program
func Fatal(msg string) {
	Logger.Fatal("[FATAL]:", msg)
}
