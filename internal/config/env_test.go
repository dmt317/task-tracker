package config

import (
	"os"
	"testing"
)

func unsetEnvVars() {
	os.Unsetenv("PORT")
	os.Unsetenv("DB_CONN")
	os.Unsetenv("IN_MEMORY")
}

type EnvVar struct {
	key   string
	value string
	set   bool
}

func getOriginalEnv(keys []string) []EnvVar {
	originalEnv := make([]EnvVar, 0, len(keys))

	for _, key := range keys {
		val, ok := os.LookupEnv(key)
		originalEnv = append(originalEnv, EnvVar{
			key:   key,
			value: val,
			set:   ok,
		})
	}

	return originalEnv
}

func restoreOriginalEnv(originalEnv []EnvVar) {
	for _, envVar := range originalEnv {
		if envVar.set {
			os.Setenv(envVar.key, envVar.value)
		} else {
			os.Unsetenv(envVar.key)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	tests := map[string]struct {
		setEnv map[string]string
		result Config
	}{
		"load config from .env file": {
			setEnv: map[string]string{
				"PORT":      "5001",
				"DB_CONN":   "user=postgres password=postgres host=localhost port=5432 dbname=tasktracker",
				"IN_MEMORY": "False",
			},
			result: Config{
				ServerPort: "5001",
				DBConn:     "user=postgres password=postgres host=localhost port=5432 dbname=tasktracker",
				InMemory:   "False",
			},
		},

		"load config with defaults": {
			setEnv: map[string]string{},
			result: Config{
				ServerPort: "8080",
				DBConn:     "user=postgres password=secret host=localhost port=5432 dbname=tasktracker",
				InMemory:   "False",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			originalEnv := getOriginalEnv([]string{"PORT", "DB_CONN", "IN_MEMORY"})
			defer restoreOriginalEnv(originalEnv)

			unsetEnvVars()

			for k, v := range test.setEnv {
				os.Setenv(k, v)
			}

			config := LoadConfig()

			if *config != test.result {
				t.Fatalf("test-case: (%q); returned [%+v]; expected [%+v]", name, config, test.result)
			}
		})
	}
}
