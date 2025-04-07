package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	originalPort, hasPort := os.LookupEnv("PORT")
	defer func() {
		if hasPort {
			os.Setenv("PORT", originalPort)
		} else {
			os.Unsetenv("PORT")
		}
	}()

	assertPort := func(expected string) {
		t.Helper()

		config := LoadConfig()

		if config.ServerPort != expected {
			t.Fatalf("Expected port %s, got %s", expected, config.ServerPort)
		}
	}

	// Тест: значение по умолчанию (переменная не установлена)
	os.Unsetenv("PORT")
	assertPort("8080")

	// Тест: установка переменной окружения
	os.Setenv("PORT", "5000")
	assertPort("5000")
}

func TestGetEnv(t *testing.T) {
	originalValue, hasValue := os.LookupEnv("TEST_KEY")
	defer func() {
		if hasValue {
			os.Setenv("TEST_KEY", originalValue)
		} else {
			os.Unsetenv("TEST_KEY")
		}
	}()

	assertEnv := func(key, defaultValue, expected string) {
		t.Helper()

		if value := getEnv(key, defaultValue); value != expected {
			t.Fatalf("getEnv(%s) = %s; want %s", key, value, expected)
		}
	}

	// Тест: если переменная НЕ установлена, должно вернуться значение по умолчанию
	os.Unsetenv("TEST_KEY")
	assertEnv("TEST_KEY", "default_value", "default_value")

	// Тест: если переменная установлена, должно вернуться её значение
	os.Setenv("TEST_KEY", "env_value")
	assertEnv("TEST_KEY", "default_value", "env_value")
}
