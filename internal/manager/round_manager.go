package manager

const (
	ticksPerRound = 5 * 60 * 3
)

type RoundManager struct {
	ticks int
}

func (r *RoundManager) Restart() {
	r.ticks = 0
}

func (r *RoundManager) Tick() {
	r.ticks++
}

func (r *RoundManager) HasEnded() bool {
	return r.ticks == ticksPerRound
}
