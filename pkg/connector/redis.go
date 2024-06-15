package connector

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type redisService struct {
  *redis.Client
}

func NewRedisService(addr, password string, db int) (*redisService, error) {
  options := &redis.Options{
    Addr: addr,
    Password: password,
    DB: db,
  }
  client := redis.NewClient(options)

  _, err := client.Ping(context.Background()).Result()
  if err != nil {
    return nil, err
  }

  return &redisService{client}, nil
}


