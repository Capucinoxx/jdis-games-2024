package manager

type Stage uint8

const (
  discoveryStage Stage = iota
  pointRushStage
)

const (
	ticksPerRound = 5 * 60 * 3
  ticksPointRushStage = 4 * 60 * 3
)

type ChangeStageHandler interface {
  ChangeStage(s Stage)
}

type RoundManager struct {
  currentStage Stage
	ticks int
  handler ChangeStageHandler
}

func NewRoundManager(handler ChangeStageHandler) *RoundManager {
  return &RoundManager{
    currentStage: discoveryStage,
    ticks: 0,
    handler: handler,
  }
}

func (r *RoundManager) Restart() {
	r.ticks = 0
  r.currentStage = discoveryStage
  r.handler.ChangeStage(discoveryStage)
}

func (r *RoundManager) Tick() {
	r.ticks++
  if r.ticks == ticksPointRushStage {
    r.handler.ChangeStage(pointRushStage)
  }
}

func (r *RoundManager) HasEnded() bool {
	return r.ticks == ticksPerRound
}
