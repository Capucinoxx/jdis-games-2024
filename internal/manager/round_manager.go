package manager

type Stage uint8

const (
  DiscoveryStage Stage = iota
  PointRushStage
)

const (
	ticksPerRound = 5 * 60 * 3
  TicksPointRushStage = 4 * 60 * 3
)


type ChangeStageHandler interface {
  ChangeStage()
}

type StageHandler struct {
  F ChangeStageHandler
  Stage Stage
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
    handler.F.ChangeStage() 
  }
}

func (r *RoundManager) Tick() {
	r.ticks++

  if handler, ok := r.handlers[r.ticks]; ok {
    handler.F.ChangeStage()
  }
}

func (r *RoundManager) AddChangeStageHandler(tick int, stage Stage, cb ChangeStageHandler) {
  r.handlers[tick] = StageHandler{ F: cb, Stage: stage  }
}

func (r *RoundManager) HasEnded() bool {
	return r.ticks == ticksPerRound
}


