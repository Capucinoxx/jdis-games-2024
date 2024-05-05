package config

import (
	"fmt"
	"os"
)

func MongoDNS() string {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"))
	return uri
}

func MongoDatabase() string {
	return os.Getenv("MONGO_DB")
}
