package env

import (
	"log"
	"os"
)

func Default(name string, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}

	return val
}

func Must(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("%s must be set", name)
	}

	return val
}
