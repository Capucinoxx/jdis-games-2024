package manager

// Package manager provides an abstraction for managing game stages and rounds
// in the application. This package includes functionalities for handling the
// state of the game, managing round ticks, and changing game stages based on
// predefined rules.

import (
	"github.com/capucinoxx/jdis-games-2024/consts"
	"github.com/capucinoxx/jdis-games-2024/pkg/model"
)

// Stage represents a game stage.
type Stage uint8

// StageHandler is an interface for handling changes in game stages.
type StageHandler interface {
	ChangeStage(state *model.GameState)
}

// RoundManager manages the ticks and state of the game rounds.
type RoundManager struct {
	ticks    int
	state    *model.GameState
	handlers map[int]StageHandler
}

// NewRoundManager creates a new instance of RoundManager.
func NewRoundManager() *RoundManager {
	return &RoundManager{
		ticks:    0,
		state:    nil,
		handlers: make(map[int]StageHandler),
	}
}

func (r *RoundManager) SetState(state *model.GameState) {
	r.state = state
}

// Restart resets the ticks and triggers the stage handler for the initial tick.
func (r *RoundManager) Restart() {
	r.ticks = 0

	if handler, ok := r.handlers[r.ticks]; ok {
		handler.ChangeStage(r.state)
	}
}

// Tick increments the tick count and triggers the stage handler for the current tick.
func (r *RoundManager) Tick() {
	r.ticks++

	if handler, ok := r.handlers[r.ticks]; ok {
		handler.ChangeStage(r.state)
	}
}

func (r *RoundManager) CurrentTick() int {
	return r.ticks / 10
}

// CurrentRound returns the current round based on the tick count.
func (r *RoundManager) CurrentRound() int8 {
	if r.ticks < consts.TicksPointRushStage {
		return 0
	}
	return 1
}

func (r *RoundManager) AddChangeStageHandler(tick int, cb StageHandler) {
	r.handlers[tick] = cb
}

func (r *RoundManager) HasEnded() bool {
	return r.ticks == consts.TicksPerRound
}

// DiscoveryStage represents the discovery stage of the game. (first stage)
type DiscoveryStage struct{}

func (s DiscoveryStage) ChangeStage(state *model.GameState) {
	spawns := state.Map.Spawns(0)
	state.SetSpawns(spawns)

	coins := make([]*model.Scorer, 0, consts.NumCoins)
	for i := 0; i < consts.NumCoins; i++ {
		coins = append(coins, model.NewCoin())
	}
	state.Reset(coins)
}

// // PointRushStage represents the point rush stage of the game.
type PointRushStage struct{}

func (s PointRushStage) ChangeStage(state *model.GameState) {
	state.SetSpawns(state.Map.Spawns(1))

	centroid := state.Map.Centroid()
	coins := []*model.Scorer{model.NewBigCoin(&centroid)}
	state.Reset(coins)
}
