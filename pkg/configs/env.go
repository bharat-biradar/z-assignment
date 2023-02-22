package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var err error

func init() {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env",err)
	}
}

func GetMongoURI() string {
	return os.Getenv("MONGODB_URI")
}
