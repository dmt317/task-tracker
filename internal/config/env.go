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
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
