package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Config func to get env value from key ---
func Config(key string) string {
	// load .env file
	err := godotenv.Load("dontsubmit.env")
	if err != nil {
		fmt.Print("Error loading .env file")
		panic(err)
	}
	return os.Getenv(key)

}
