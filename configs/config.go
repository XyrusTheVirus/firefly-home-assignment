package configs

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envLoaded   bool = false
	configMutex sync.Mutex
)

// loadEnv loads the .env file
func loadEnv() {
	defer configMutex.Unlock()
	envLoaded = true

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Env returns the value of the key in the .env file
func Env(key string, fallback string) string {
	if !envLoaded {
		configMutex.Lock()
		if !envLoaded {
			loadEnv()
		}
	}

	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// EnvInt returns the value of the key in the .env file as an int
func EnvInt(key string, fallback string) int {
	value, _ := strconv.Atoi(Env(key, fallback))
	return value
}
