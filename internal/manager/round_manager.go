package manager

import (
	"github.com/capucinoxx/forlorn/pkg/config"
	"github.com/capucinoxx/forlorn/pkg/model"
	"github.com/capucinoxx/forlorn/pkg/utils"
)

type Stage uint8

const (
  discoveryStage Stage = iota
  pointRushStage
)

const (
	ticksPerRound = 5 * 60 * 3
  TicksPointRushStage = 4 * 60 * 3
)




type StageHandler interface {
  ChangeStage(state *model.GameState)
}


type RoundManager struct {
	ticks int
  state *model.GameState
  handlers map[int]StageHandler
}


func NewRoundManager() *RoundManager {
  return &RoundManager{
    ticks: 0,
    state: nil,
    handlers: make(map[int]StageHandler),
  }
}

func (r *RoundManager) SetState(state *model.GameState) {
  r.state = state
}


func (r *RoundManager) Restart() {
	r.ticks = 0
  
  if handler, ok := r.handlers[r.ticks]; ok {
    handler.ChangeStage(r.state) 
  }
}


func (r *RoundManager) Tick() {
	r.ticks++

  if handler, ok := r.handlers[r.ticks]; ok {
    handler.ChangeStage(r.state)
  }
}

func (r *RoundManager) CurrentTick() int {
  return r.ticks / 10
}

func (r *RoundManager) CurrentRound() int8 {
  if r.ticks < TicksPointRushStage {
    return 0
  }
  return 1
}

func (r *RoundManager) AddChangeStageHandler(tick int, cb StageHandler) {
  r.handlers[tick] = cb
}


func (r *RoundManager) HasEnded() bool {
	return r.ticks == ticksPerRound
}


type DiscoveryStage struct {}
func (s DiscoveryStage) ChangeStage(state *model.GameState) {
  spawns := state.Map.Spawns(0)
  state.SetSpawns(spawns)
  utils.Log("stage", "stage", "DISCOVERY STAGE CHANGE STAGE")

  coins := make([]*model.Scorer, 0, config.NumCoins)
  for i := 0; i < config.NumCoins; i++ {
    coins = append(coins, model.NewCoin())
  }
  state.Reset(coins)
}


// TODO: reset players avec nouvelle pos (full health, 0 bullets)
// TODO: effacer toute les pièces
// TODO: ajouter pièces au milieu
type PointRushStage struct {}
func (s PointRushStage) ChangeStage(state *model.GameState) {
  state.SetSpawns(state.Map.Spawns(1))
  utils.Log("stage", "stage", "POINT RUSH STAGE CHANGE STAGE")
  
  centroid := state.Map.Centroid()
  coins := []*model.Scorer{model.NewBigCoin(&centroid)}
  state.Reset(coins)
}
