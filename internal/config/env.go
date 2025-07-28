package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING:", err)
	}

	return &Config{
		ServerPort: getEnv("PORT", "8080"),
		DBConn:     getEnv("DB_CONN", "user=postgres password=secret host=localhost port=5432 dbname=tasktracker"),
		AuthDBConn: getEnv("AUTH_DB_CONN", "user=postgres password=secret host=localhost port=5432 dbname=auth"),
		InMemory:   getEnv("IN_MEMORY", "False"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
