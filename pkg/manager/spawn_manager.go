package manager

import (
	"math/rand"
	"time"

	"github.com/capucinoxx/forlorn/pkg/config"
	"github.com/capucinoxx/forlorn/pkg/model"
  "github.com/capucinoxx/forlorn/pkg/utils"
)

type SpawnManager interface {
  Spawn() *model.Point
}

type RandomSpawnManager struct {
  r *rand.Rand
}


func NewRandomSpawnManager() SpawnManager {
  return &RandomSpawnManager{
    r: rand.New(rand.NewSource(time.Now().UnixNano())),
  }
}


func (s *RandomSpawnManager) Spawn() *model.Point {
  x := s.r.Float64() * float64(config.MapWidth * config.CellWidth)
  y := s.r.Float64() * float64(config.MapWidth * config.CellWidth)
  return &model.Point{X: utils.Round(x, 2), Y: utils.Round(y, 2)}
}


type ListSpawnManager struct {
  spawns []*model.Point
  index int
}


func NewListSpawnManager(spawns []*model.Point) SpawnManager {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  utils.Shuffle(r, spawns)
  return &ListSpawnManager{
    spawns: spawns,
    index: 0,
  }
}


func (s *ListSpawnManager) Spawn() *model.Point {
  spawn := s.spawns[s.index]
  s.index = (s.index + 1) % len(s.spawns)
  return spawn
}
