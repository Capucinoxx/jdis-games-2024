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
	"os"
	"sync"
	"time"

	"github.com/capucinoxx/jdis-games-2024/pkg/connector"
	"github.com/capucinoxx/jdis-games-2024/pkg/model"
	"github.com/capucinoxx/jdis-games-2024/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerEntry struct {
	Score int32 `json:"y"`
	Time  int64 `json:"x"`
}

// PlayerScore represents the score details of a player.
type PlayerScore struct {
	Name    string `json:"name"`
	Score   int    `json:"score"`
	Color   int    `json:"color"`
	Ranking int    `json:"ranking"`
}

type Cache struct {
	leaderboard []PlayerScore
	histories   map[string][]PlayerEntry
	lastUpdate  time.Time
	duration    time.Duration
	mu          sync.Mutex
}

func NewCache(duration time.Duration) *Cache {
	return &Cache{
		duration:  duration,
		histories: make(map[string][]PlayerEntry),
	}
}

func (c *Cache) Get() ([]PlayerScore, map[string][]PlayerEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Since(c.lastUpdate) < c.duration {
		return c.leaderboard, c.histories, true
	}
	return nil, nil, false
}

func (c *Cache) Set(leaderboard []PlayerScore, histories map[string][]PlayerEntry) {
	c.leaderboard = leaderboard
	c.histories = histories
	c.lastUpdate = time.Now()
}

// ScoreManager handles the management of player scores.
type ScoreManager struct {
	redis        *connector.RedisService
	mongo        *connector.MongoService
	currentScore map[string]int
	mu           sync.Mutex
	visible      bool
	cache        *Cache
	persist      bool
}

// NewScoreManager creates a new ScoreManager with the specified Redis and MongoDB services.
func NewScoreManager(redis *connector.RedisService, mongo *connector.MongoService) *ScoreManager {
	return &ScoreManager{
		redis:        redis,
		mongo:        mongo,
		currentScore: make(map[string]int),
		visible:      true,
		cache:        NewCache(time.Minute),
		persist:      os.Getenv("RANK") == "RANKED",
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
	if !sm.persist {
		return nil
	}

	val, err := sm.redis.ZRevRangeWithScores(context.Background(), "leaderboard", 0, -1).Result()
	if err != nil {
		return err
	}
	now := time.Now()

	var errs utils.Errors
	for _, s := range val {
		if err = sm.mongo.Push("scores", s.Member.(string), "scores", bson.M{"score": int32(s.Score), "time": now}); err != nil {
			errs.Append(err)
		}
	}
	return errs.Error()
}

// Add increments the score of a player identified by UUID in the Redis leaderboard.
func (sm *ScoreManager) Adds(players []model.PlayerScore) {
	if !sm.persist {
		return
	}

	go func() {
		ctx := context.Background()
		pipe := sm.redis.Pipeline()

		for _, player := range players {
			pipe.ZIncrBy(ctx, "leaderboard", float64(player.Score), player.Name)
		}

		_, err := pipe.Exec(ctx)
		if err != nil {
			utils.Log("error", "redis", "adds players scores %s", err)
		}
	}()
}

func (sm *ScoreManager) findColors(names []string) ([]int, error) {
	filter := bson.M{"username": bson.M{"$in": names}}
	res, err := sm.mongo.Find("users", filter)
	if err != nil {
		return nil, err
	}

	colors := make([]int, len(res))
	for _, v := range res {
		for j, name := range names {
			if name == v["username"].(string) {
				colors[j] = int(v["color"].(int32))
				break
			}
		}
	}
	return colors, nil
}

// Rank retrieves the ranked player scores from the Redis leaderboard.
// It returns a map of player UUIDs to their respective scores and positions.
func (sm *ScoreManager) Rank() ([]PlayerScore, map[string][]PlayerEntry, error) {
	if !sm.persist {
		return []PlayerScore{}, map[string][]PlayerEntry{}, nil
	}

	if leaderboard, histories, ok := sm.cache.Get(); ok {
		return leaderboard, histories, nil
	}
	sm.cache.mu.Lock()
	defer sm.cache.mu.Unlock()

	val, err := sm.redis.ZRevRangeWithScores(context.Background(), "leaderboard", 0, -1).Result()
	if err != nil {
		return nil, nil, err
	}

	leaderboard := make([]PlayerScore, len(val))
	names := make([]string, 10)

	for i, v := range val {
		uuid := v.Member.(string)
		leaderboard[i] = PlayerScore{
			Ranking: i + 1,
			Name:    uuid,
			Score:   int(v.Score),
		}

		if i < 10 {
			names[i] = uuid
		}
	}
	colors, err := sm.findColors(names)
	if err != nil {
		return nil, nil, err
	}

	for i := range leaderboard {
		leaderboard[i].Color = colors[i]
	}

	filter := bson.M{"_id": bson.M{"$in": names}}
	res, err := sm.mongo.Find("scores", filter)
	if err != nil {
		return nil, nil, err
	}

	histories := make(map[string][]PlayerEntry)
	for _, history := range res {
		name := history["_id"].(string)
		entries, ok := history["scores"].(primitive.A)
		if !ok {
			utils.Log("error", "persist", "error retrieve scores")
			continue
		}

		histories[name] = make([]PlayerEntry, 0, len(entries))
		for _, entry := range entries {
			doc, ok := entry.(bson.M)
			if !ok {
				utils.Log("error", "persist", "error casting score to bson.M")
				continue
			}

			score, ok := doc["score"].(int32)
			if !ok {
				utils.Log("error", "persist", "error retrieving score as int %v", doc["score"])
				continue
			}

			scoreTime, ok := doc["time"].(primitive.DateTime)
			if !ok {
				utils.Log("error", "persist", "error retrieving score time as time.Time, %v", doc["time"])
				continue
			}

			histories[name] = append(histories[name], PlayerEntry{Score: score, Time: int64(scoreTime.Time().UnixMilli() / 10_000)})
		}
	}
	sm.cache.Set(leaderboard, histories)

	return leaderboard, histories, nil
}
