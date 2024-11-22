package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("LOG_EXECUTION_ID") == "" {
		err := godotenv.Load()
		if err != nil {
			panic(fmt.Sprintf("Error loading .env file: %v", err))
		}
	}
}

func Getenv(key string) string {
	return os.Getenv(key)
}
