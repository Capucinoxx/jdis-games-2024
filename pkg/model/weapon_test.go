package model

import (
	"math"
	"testing"

	"github.com/capucinoxx/forlorn/consts"
)

func TestBladeIntersection(t *testing.T) {

	tests := map[string]struct {
		blade          *Blade
		otherPlayers   []*Player
		rotation       float64
		expectedHealth []int
	}{
		"Blade intersects stationary player": {
			blade:          NewBlade(NewPlayer("0", 0, &Point{X: 0, Y: 0}, nil)),
			otherPlayers:   []*Player{NewPlayer("1", 0, &Point{X: 1, Y: 0}, nil)},
			rotation:       0.0,
			expectedHealth: []int{100 - consts.BladeDmg},
		},
		"Blade rotates and intersects player": {
			blade:          NewBlade(NewPlayer("0", 0, &Point{X: 0, Y: 0}, nil)),
			otherPlayers:   []*Player{NewPlayer("1", 0, &Point{X: 0, Y: 1}, nil)},
			rotation:       math.Pi / 2,
			expectedHealth: []int{100 - consts.BladeDmg},
		},
		"Blade does not intersect player": {
			blade:          NewBlade(NewPlayer("0", 0, &Point{X: 0, Y: 0}, nil)),
			otherPlayers:   []*Player{NewPlayer("1", 0, &Point{X: 1, Y: 0}, nil)},
			rotation:       math.Pi / 2,
			expectedHealth: []int{100},
		},
		"Blade intersects two players": {
			blade:          NewBlade(NewPlayer("0", 0, &Point{X: 0, Y: 0}, nil)),
			otherPlayers:   []*Player{NewPlayer("1", 0, &Point{X: 0, Y: 1}, nil), NewPlayer("2", 0, &Point{X: 0, Y: 1}, nil)},
			rotation:       math.Pi / 2,
			expectedHealth: []int{100 - consts.BladeDmg, 100 - consts.BladeDmg},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.blade.Update(tt.otherPlayers, nil, &tt.rotation)
			for i := range tt.otherPlayers {
				if tt.otherPlayers[i].health != tt.expectedHealth[i] {
					t.Errorf("Player[%d] health mismatch: got %d, want %d", i, tt.otherPlayers[i].health, tt.expectedHealth[i])
				}
			}
		})
	}
}
