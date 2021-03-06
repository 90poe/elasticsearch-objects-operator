package utils

import (
	"log"
	"os"
)

//MustGetEnv would get Env variable of die
func MustGetEnv(name string) string {
	value, found := os.LookupEnv(name)
	if !found {
		log.Fatalf("environment variable %s is missing", name)
	}
	return value
}

// GetEnv would try to get ENV variable
func GetEnv(name string) string {
	return os.Getenv(name)
}
