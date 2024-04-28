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

	config.VerifyEndpoint = os.Getenv("VERIFY_ENDPOINT")

	config.ChannelBaseURL = os.Getenv("CHANNEL_BASE_URL")

	config.ResourceServerPort = os.Getenv("RESOURCE_SERVER_PORT")
}
