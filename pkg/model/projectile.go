package model

import (
	"math"

	"github.com/capucinoxx/forlorn/pkg/config"
)

// Projectile represents a moving projectile in the game.
type Projectile struct {
  Object
  Destination *Point
}

func NewProjectile(pos *Point, dest *Point) *Projectile {
  p := &Projectile{Destination: dest}
  p.setup(pos, config.ProjectileSize)

  return p
} 

// ApplyMovement updates the projectile's position based on its direction and a delta time.
func (p *Projectile) moveToDestination(dt float32) bool {
  dest := p.Destination

  dx := float64(dest.X - p.Position.X)
  dy := float64(dest.Y - p.Position.Y)
  dist := math.Abs(dx) + math.Abs(dy)

  if dist > config.ProjectileSpeed * float64(dt) {
    nextX := p.Position.X + float32(dx/dist * config.ProjectileSpeed) * dt
    nextY := p.Position.Y + float32(dy/dist * config.ProjectileSpeed) * dt

    p.Position.X = nextX
    p.Position.Y = nextY
    return true
  } else {
    p.Remove()
    return false
  } 
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
		if !p.moveToDestination(dt) {
      continue
    }

    if p.Position.Equals(p.Destination, 0.2) {
      p.Remove()
      continue
    } 

		for _, enemy := range players {
			if c.owner.Nickname == enemy.Nickname {
				continue
			}

			if p.IsCollidingWithPlayer(enemy) {
				enemy.TakeDmg(config.ProjectileDmg)
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
  c.Projectiles = append(c.Projectiles, NewProjectile(
    &Point{X: c.owner.Collider.Pivot.X, Y: c.owner.Collider.Pivot.Y},
    &Point{X: pos.X, Y: pos.Y},
  ))
}
