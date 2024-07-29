package model

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/capucinoxx/jdis-games-2024/consts"
)

func TestCannonShootAt(t *testing.T) {
	owner := NewPlayer("owner", 0, &Point{X: 5, Y: 5}, nil)

	tests := map[string]struct {
		projectileTarget Point
		expectedPosition Point
	}{
		"Shoot positive destination": {
			projectileTarget: Point{X: 10, Y: 10},
			expectedPosition: Point{
				X: 5.0 + consts.ProjectileSpeed/math.Sqrt(2),
				Y: 5.0 + consts.ProjectileSpeed/math.Sqrt(2),
			},
		},
		"Shoot negative destination": {
			projectileTarget: Point{X: -10, Y: -10},
			expectedPosition: Point{
				X: 5.0 - consts.ProjectileSpeed/math.Sqrt(2),
				Y: 5.0 - consts.ProjectileSpeed/math.Sqrt(2),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cannon := NewCanon(owner)

			cannon.ShootAt(tt.projectileTarget)

			cannon.Update([]*Player{}, 1.0)

			pos := cannon.Projectiles[0].Position
			expected := tt.expectedPosition
			if math.Abs(pos.X-expected.X) > 0.0001 || math.Abs(pos.Y-expected.Y) > 0.0001 {
				t.Errorf("Projectile position (%v, %v) != expected position (%v, %v)", pos.X, pos.Y, expected.X, expected.Y)
			}
		})
	}
}

func TestCannonCollisionCases(t *testing.T) {
	owner := NewPlayer("owner", 0, &Point{X: 0, Y: 0}, nil)

	tests := map[string]struct {
		enemyPositions   []*Point
		expectedHealths  []int
		projectileTarget Point
	}{
		"Single enemy hit": {
			enemyPositions:   []*Point{{X: 5, Y: 5}},
			expectedHealths:  []int{100 - consts.ProjectileDmg},
			projectileTarget: Point{X: 5, Y: 5},
		},
		"Multiple enemies hit": {
			enemyPositions:   []*Point{{X: 5, Y: 5}, {X: 6, Y: 5}},
			expectedHealths:  []int{100 - consts.ProjectileDmg, 100},
			projectileTarget: Point{X: 5, Y: 5},
		},
		"No enemies hit": {
			enemyPositions:   []*Point{{X: 15, Y: 15}, {X: 16, Y: 16}},
			expectedHealths:  []int{100, 100},
			projectileTarget: Point{X: 5, Y: 5},
		},
		"Enemy outside of range": {
			enemyPositions:   []*Point{{X: 100, Y: 100}},
			expectedHealths:  []int{100},
			projectileTarget: Point{X: 5, Y: 5},
		},
		"Projectile timed out before reach cible": {
			enemyPositions:   []*Point{{X: consts.ProjectileSpeed * (consts.ProjectileTTL + 0.3), Y: consts.ProjectileSpeed * (consts.ProjectileTTL + 0.3)}},
			expectedHealths:  []int{100},
			projectileTarget: Point{X: 1000, Y: 1000},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cannon := NewCanon(owner)
			enemies := make([]*Player, len(tt.enemyPositions))
			for i, pos := range tt.enemyPositions {
				enemies[i] = NewPlayer("enemy"+fmt.Sprint(i), 100, pos, nil)
			}

			cannon.ShootAt(tt.projectileTarget)

			for dt := 0.3; dt < 2.0; dt += 0.3 {
				cannon.Update(enemies, dt)
			}

			for i, enemy := range enemies {
				if enemy.health != tt.expectedHealths[i] {
					t.Errorf("Enemy[%d] health incorrect after being hit, got %d, want %d", i, enemy.health, tt.expectedHealths[i])
				}
			}

			for _, p := range cannon.Projectiles {
				if p.IsAlive() {
					t.Errorf("Expected projectile to be removed, but it is still alive")
				}
			}
		})
	}
}
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
			tt.blade.Update(tt.otherPlayers, &tt.rotation)
			for i := range tt.otherPlayers {
				if tt.otherPlayers[i].health != tt.expectedHealth[i] {
					t.Errorf("Player[%d] health mismatch: got %d, want %d", i, tt.otherPlayers[i].health, tt.expectedHealth[i])
				}
			}
		})
	}
}

func TestBladeIntersectoionFuzzing(t *testing.T) {
	const numTests = 10_000

	const maxDistance = 10.0

	for i := 0; i < numTests; i++ {
		t.Run(fmt.Sprintf("FuzzTest-%d", i), func(t *testing.T) {
			ownerPos := &Point{X: rand.Float64() * maxDistance, Y: rand.Float64() * maxDistance}
			enemyPos := &Point{X: rand.Float64() * maxDistance, Y: rand.Float64() * maxDistance}

			owner := NewPlayer("owner", 0, ownerPos, nil)
			enemy := NewPlayer("enemy", 0, enemyPos, nil)
			blade := NewBlade(owner)

			distance := math.Sqrt(math.Pow(ownerPos.X-enemyPos.X, 2) + math.Pow(ownerPos.Y-enemyPos.Y, 2))

			rotation := rand.Float64() * 2 * math.Pi
			rotation = math.Pi / 4.0
			blade.Update([]*Player{enemy}, &rotation)

			if distance > ((consts.PlayerSize/2.0 + consts.BladeSize) + math.Cos(math.Pi/4.0)) {
				if enemy.health != 100 {
					t.Errorf("Enemy should not take damage when distance > %f, but got health %d", ((consts.PlayerSize/2.0 + consts.BladeSize) + math.Cos(math.Pi/4.0)), enemy.health)
				}
			}
		})
	}
}
