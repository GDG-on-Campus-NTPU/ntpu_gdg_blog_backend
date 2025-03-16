package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Fprintln(os.Stderr, "\033[31mError loading .env file If running in cloud run ignore this error\033[0m")
	}
}

func Getenv(key string) string {
	return os.Getenv(key)
}
