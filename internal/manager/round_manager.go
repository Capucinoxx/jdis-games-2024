package manager

import "github.com/capucinoxx/forlorn/pkg/utils"

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
  ChangeStage()
}

type RoundManager struct {
	ticks int
  handlers map[int]StageHandler
}

func NewRoundManager() *RoundManager {
  return &RoundManager{
    ticks: 0,
    handlers: make(map[int]StageHandler),
  }
}

func (r *RoundManager) Restart() {
	r.ticks = 0
  
  if handler, ok := r.handlers[r.ticks]; ok {
    handler.ChangeStage() 
  }
}

func (r *RoundManager) Tick() {
	r.ticks++

  if handler, ok := r.handlers[r.ticks]; ok {
    handler.ChangeStage()
  }
}

func (r *RoundManager) AddChangeStageHandler(tick int, cb StageHandler) {
  r.handlers[tick] = cb
}

func (r *RoundManager) HasEnded() bool {
	return r.ticks == ticksPerRound
}

type DiscoveryStage struct {}
func (s DiscoveryStage) ChangeStage() {
  utils.Log("game", "stage", "Set DiscoveryStage")
}

// TODO: reset players avec nouvelle pos (full health, 0 bullets)
// TODO: effacer toute les pièces
// TODO: ajouter pièces au milieu
type PointRushStage struct {}
func (s PointRushStage) ChangeStage() {
  utils.Log("game", "stage", "Set PointRushStage")
}
