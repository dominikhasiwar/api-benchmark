package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, defaultValue string) string {
	// Load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
