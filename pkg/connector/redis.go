package connector

// RedisService provides functionality for communicating with a Redis database.
// This service allows for connecting to the Redis server, performing common Redis operations,
// and managing the connection lifecycle.
//
// This service ensures efficient resource management by properly handling
// database connections and operations.
//
// Usage of this package involves creating an instance of RedisService with the Redis server address,
// password, and database index. The RedisService instance can then be used to perform various
// Redis operations.

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisService provides functionality for communicating with a Redis database.
type RedisService struct {
	*redis.Client
}

// NewRedisService creates a new Redis service. Upon creation, it establishes a connection
// to the Redis server. If the connection fails, an error is returned.
func NewRedisService(addr, password string, db int) (*RedisService, error) {
	options := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	client := redis.NewClient(options)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisService{client}, nil
}
