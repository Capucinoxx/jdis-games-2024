package manager

import (
	"context"
	"sync"
	"time"

	"github.com/capucinoxx/forlorn/pkg/connector"
	"github.com/capucinoxx/forlorn/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type PlayerScores = map[string]PlayerScore 

type PlayerScore struct {
  Position int `json:"pos"`
  UUID string `json:"uuid"`
  Score float64 `json:"score"`
}

type ScoreManager struct {
  redis *connector.RedisService
  mongo *connector.MongoService
  currentScore map[string]int
  mu sync.Mutex
}

func NewScoreManager(redis *connector.RedisService, mongo *connector.MongoService) *ScoreManager {
  return &ScoreManager{
    redis: redis,
    mongo: mongo,
    currentScore: make(map[string]int),
  }
}

func (sm *ScoreManager) Persist() error {
  scores, err := sm.Rank()
  if err != nil {
    return err
  }
  now := time.Now()

  var errs utils.Errors 
  for _, s := range scores {
    if err = sm.mongo.Push("users", s.UUID, "scores", bson.M{ "score": s.Score, "time": now }); err != nil {
      errs.Append(err)
    }
  }
  return errs.Error()
}

func (sm *ScoreManager) Add(uuid string, score float64) error {
  return sm.redis.ZIncrBy(context.Background(), "leaderboard", score, uuid).Err()
}

func (sm *ScoreManager) Rank() (PlayerScores, error) {
  val, err := sm.redis.ZRevRangeWithScores(context.Background(), "leaderboard", 0, -1).Result()
  if err != nil {
    return PlayerScores{}, err 
  }

  players := make(PlayerScores)
  for i, v := range val {
    uuid := v.Member.(string)
    players[uuid] = PlayerScore{
      Position: i,
      UUID: uuid,
      Score: v.Score,
    }
  }

  return players, err  
}

