package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/capucinoxx/forlorn/pkg/manager"
	"github.com/capucinoxx/forlorn/pkg/utils"
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

func RedisAddr() string {
	return os.Getenv("REDIS_ADDR")
}

func RedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func RequiredAdmins() []manager.TokenInfo {
	var admins []manager.TokenInfo
	_ = json.Unmarshal([]byte(os.Getenv("ADMINS")), &admins)
	utils.Log("config", "admins", "%d admins have been retrieved", len(admins))
	return admins
}

func Port() int {
	v, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return 8087
	}
	return v
}
