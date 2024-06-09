package model

import (
	"math"

	"github.com/google/uuid"
)

const PROJECTILE_DMG = 30

const (
  projectile_speed float64 = 1
)

// Projectile represents a moving projectile in the game.
type Projectile struct {
  uuid [16]byte
	Position  *Point
	Destination *Point

	// // cleanup indicates whether the projectile should be removed from the game.
	cleanup bool
}

// ApplyMovement updates the projectile's position based on its direction and a delta time.
func (p *Projectile) moveToDestination(dt float32) {
  dest := p.Destination

  dx := float64(dest.X - p.Position.X)
  dy := float64(dest.Y - p.Position.Y)
  dist := math.Abs(dx) + math.Abs(dy)

  if dist > projectile_speed * float64(dt) {
    nextX := p.Position.X + float32(dx/dist * projectile_speed) * dt
    nextY := p.Position.Y + float32(dy/dist * projectile_speed) * dt

    p.Position.X = nextX
    p.Position.Y = nextY
  } 
}

// IsCollidingWithPlayer checks if the projectile is colliding with a given player.
func (p *Projectile) IsCollidingWithPlayer(player *Player) bool {
	rect := player.Collider.rect
	polygon := []*Point{rect.a, rect.b, rect.c, rect.d}
	return p.Position.IsInPolygon(polygon)
}

// IsCollidingWithEnvironment checks if the projectile is colliding with any non-projectile colliders in the map.
func (p *Projectile) IsCollidingWithEnvironment(m Map) bool {
	for _, collider := range m.Colliders() {
		if collider.Type == ColliderProjectile {
			continue
		}

		if p.Position.IsInPolygon(collider.Points) {
			return true
		}
	}

	return false
}

// Remove marks the projectile for cleanup.
func (p *Projectile) Remove() { p.cleanup = true }

// IsAlive indicating if the projectile is still active.
func (p *Projectile) IsAlive() bool { return !p.cleanup }

// Cannon represents a cannon that can shoot projectiles.
type Cannon struct {
	Projectiles []*Projectile
	owner       *Player
}

// NewCanon creates a new cannon associated with a specific player.
func NewCanon(owner *Player) *Cannon {
	return &Cannon{
		owner: owner,
	}
}

// Update processes all projectiles for movement and collision detection.
func (c *Cannon) Update(players []*Player, m Map, dt float32) {
	for _, p := range c.Projectiles {
		p.moveToDestination(dt)

    if p.Position.Equals(p.Destination, 0.2) {
      p.Remove()
      continue
    } 

		for _, enemy := range players {
			if c.owner.ID == enemy.ID {
				continue
			}

			if p.IsCollidingWithPlayer(enemy) {
				enemy.TakeDmg(PROJECTILE_DMG)
        p.Remove()
				continue
			}

			if p.IsCollidingWithEnvironment(m) {
				p.Remove()
			}
		}
	}

	// Filters out projectiles that need to be cleaned up.
	projectiles := make([]*Projectile, 0)
	for _, p := range c.Projectiles {
    if p.IsAlive() {
      projectiles = append(projectiles, p)
    }
  }
  c.Projectiles = projectiles
}

// ShootAt creates a projectile at a specified position and calculates its direction.
func (c *Cannon) ShootAt(pos Point) {
  id := uuid.New()
  projectile := &Projectile{
    Position: &Point{X: c.owner.Collider.Pivot.X, Y: c.owner.Collider.Pivot.Y},
    Destination: &Point{X: pos.X, Y: pos.Y},
    cleanup: false,
  }

  copy(projectile.uuid[:], id[:])
  c.Projectiles = append(c.Projectiles, projectile)
}
