package manager

import "github.com/capucinoxx/forlorn/pkg/connector"

type ScoreManager struct {
  redis *connector.RedisService
  mongo *connector.MongoService
}

func NewScoreManager(redis *connector.RedisService, mongo *connector.MongoService) *ScoreManager {
  return &ScoreManager{
    redis: redis,
    mongo: mongo,
  }
}
