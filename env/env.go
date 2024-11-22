package env

import (
	"fmt"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("LOG_EXECUTION_ID") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			runtime.Goexit()
		}
	}
}

func Getenv(key string) string {
	return os.Getenv(key)
}
