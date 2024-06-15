package connector

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
  *redis.Client
}

func NewRedisService(addr, password string, db int) (*RedisService, error) {
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

  return &RedisService{client}, nil
}


