package manager

// ScoreManager handles the management of player scores using Redis and MongoDB services.
// This manager provides functionality for toggling score visibility, persisting scores,
// adding new scores, and retrieving ranked player scores.
//
// This manager uses the RedisService for real-time score updates and the MongoService
// for persistent storage of historical scores. It also ensures thread safety using a mutex
// for concurrent access to visibility status.

import (
	"context"
	"sync"
	"time"

	"github.com/capucinoxx/forlorn/pkg/connector"
	"github.com/capucinoxx/forlorn/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type PlayerScores = map[string]PlayerScore

// PlayerScore represents the score details of a player.
type PlayerScore struct {
	Position int     `json:"pos"`
	UUID     string  `json:"uuid"`
	Score    float64 `json:"score"`
}

// ScoreManager handles the management of player scores.
type ScoreManager struct {
	redis        *connector.RedisService
	mongo        *connector.MongoService
	currentScore map[string]int
	mu           sync.Mutex
	visible      bool
}

// NewScoreManager creates a new ScoreManager with the specified Redis and MongoDB services.
func NewScoreManager(redis *connector.RedisService, mongo *connector.MongoService) *ScoreManager {
	return &ScoreManager{
		redis:        redis,
		mongo:        mongo,
		currentScore: make(map[string]int),
		visible:      true,
	}
}

// ToggleVisibility toggles the visibility of scores and returns the new visibility status.
func (sm *ScoreManager) ToggleVisibility() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.visible = !sm.visible

	return sm.visible
}

// IsVisible returns the current visibility status of scores.
func (sm *ScoreManager) IsVisible() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.visible
}

// Persist saves the current scores to MongoDB. It retrieves the ranked scores from Redis,
// associates them with the current time, and pushes them to the MongoDB collection.
func (sm *ScoreManager) Persist() error {
	scores, err := sm.Rank()
	if err != nil {
		return err
	}
	now := time.Now()

	var errs utils.Errors
	for _, s := range scores {
		if err = sm.mongo.Push("users", s.UUID, "scores", bson.M{"score": s.Score, "time": now}); err != nil {
			errs.Append(err)
		}
	}
	return errs.Error()
}

// Add increments the score of a player identified by UUID in the Redis leaderboard.
func (sm *ScoreManager) Add(uuid string, score float64) error {
	return sm.redis.ZIncrBy(context.Background(), "leaderboard", score, uuid).Err()
}

// Rank retrieves the ranked player scores from the Redis leaderboard.
// It returns a map of player UUIDs to their respective scores and positions.
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
			UUID:     uuid,
			Score:    v.Score,
		}
	}

	return players, err
}
