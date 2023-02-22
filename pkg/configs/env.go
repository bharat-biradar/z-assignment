package configs

import (
	"os"
)

func GetMongoURI() string {
	return os.Getenv("MONGODB_URI")
}

func GetPort() string {
	return os.Getenv("PORT")
}
