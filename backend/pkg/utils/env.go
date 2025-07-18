package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadEnvFile() {
	env := GetEnv("GO_ENV", "development")
	if env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found or could not be loaded: %v", err)
		} else {
			log.Println("Loaded .env file for development environment")
		}
	}
}