package model

import (
	"fmt"
	"testing"

	"github.com/capucinoxx/jdis-games-2024/consts"
)

func TestScorerCollision(t *testing.T) {
	tests := map[string]struct {
		playerPositions []*Point
		expectedScores  []int
		scorerPositions []*Point
		expectedChanges []bool
	}{
		"Coin collision": {
			playerPositions: []*Point{{X: 5, Y: 5}},
			expectedScores:  []int{int(consts.CoinValue)},
			scorerPositions: []*Point{{X: 5, Y: 5}},
			expectedChanges: []bool{true},
		},
		"No collision": {
			playerPositions: []*Point{{X: 0, Y: 0}},
			expectedScores:  []int{0},
			scorerPositions: []*Point{{X: 10, Y: 10}},
			expectedChanges: []bool{false},
		},
		"Multiple players and coins collision": {
			playerPositions: []*Point{{X: 5, Y: 5}, {X: 10, Y: 10}},
			expectedScores:  []int{int(consts.CoinValue), int(consts.CoinValue)},
			scorerPositions: []*Point{{X: 5, Y: 5}, {X: 10, Y: 10}, {X: 1_000, Y: 1_000}, {X: 1_000, Y: 1_000}},
			expectedChanges: []bool{true, true},
		},
		"Partial collision": {
			playerPositions: []*Point{{X: 5, Y: 5}, {X: 15, Y: 15}},
			expectedScores:  []int{int(consts.CoinValue), 0},
			scorerPositions: []*Point{{X: 5, Y: 5}, {X: 10, Y: 10}, {X: 1_000, Y: 1_000}},
			expectedChanges: []bool{true, false},
		},
		"two player take in same time coin": {
			playerPositions: []*Point{{X: 5, Y: 5}, {X: 5, Y: 5}},
			expectedScores:  []int{int(consts.CoinValue), 0},
			scorerPositions: []*Point{{X: 5, Y: 5}, {X: 1_000, Y: 1_000}},
			expectedChanges: []bool{true},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			players := make([]*Player, len(tt.playerPositions))
			for i, pos := range tt.playerPositions {
				players[i] = NewPlayer(fmt.Sprintf("Player%d", i), 0, pos, nil)
			}

			scorers := NewScorers()
			initialUUIDs := make([][16]byte, 0, len(tt.scorerPositions))
			for _, pos := range tt.scorerPositions {
				coin := NewCoin(pos)
				scorers.Add(coin)
				initialUUIDs = append(initialUUIDs, coin.uuid)
			}

			gameState := NewGameState(nil)
			gameState.coins = scorers

			for _, player := range players {
				player.Update([]*Player{}, gameState, 0)
			}
			scorers.Update()

			for i, scorer := range scorers.List() {
				hasChanged := string(scorer.uuid[:]) != string(initialUUIDs[i][:])
				if hasChanged != tt.expectedChanges[i] {
					t.Errorf("Expected change for scorer %d to be %v, got %v", i, tt.expectedChanges[i], hasChanged)
				}
				if i == len(tt.expectedChanges)-1 {
					break
				}
			}

			for i, player := range players {
				if player.score != tt.expectedScores[i] {
					t.Errorf("Expected score for player %d to be %d, got %d", i, tt.expectedScores[i], player.score)
				}
			}
		})
	}
}
