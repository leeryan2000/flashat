package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	GO_ENV string
	// database fields
	DB_USER string
	DB_PASS string
	DB_NAME string
	DB_HOST string
	DB_PORT string

	REDIS_ADDR string
}

func LoadConfig() Configuration {
	envFile := ".env"

	// Load the environment file
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Warning: Error loading %s file, reading from system env only", envFile)
	}

	configuration := Configuration{}
	configuration.GO_ENV = os.Getenv("GO_ENV")
	configuration.DB_USER = os.Getenv("DB_USER")
	configuration.DB_PASS = os.Getenv("DB_PASS")
	configuration.DB_NAME = os.Getenv("DB_NAME")
	configuration.DB_HOST = os.Getenv("DB_HOST")
	configuration.DB_PORT = os.Getenv("DB_PORT")
	configuration.REDIS_ADDR = os.Getenv("REDIS_ADDR")

	return configuration
}
