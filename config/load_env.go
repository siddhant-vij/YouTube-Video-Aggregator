package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(config *ApiConfig) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.DatabaseURL = os.Getenv("DATABASE_URL")

	config.AuthServerPort = os.Getenv("AUTH_SERVER_PORT")
	config.AuthVerifyEndpoint = os.Getenv("AUTH_VERIFY_ENDPOINT")

	config.ResourceServerPort = os.Getenv("RESOURCE_SERVER_PORT")
}
