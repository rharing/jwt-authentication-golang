package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Config func to get env value from key ---
/*
note the dontsubmit.env file should be in the root dir
*/
func Config(key string) string {
	// load .env file
	err := godotenv.Load("dontsubmit.env")
	if err != nil {
		fmt.Print("Error loading .env file")
		panic(err)
	}
	return os.Getenv(key)

}
