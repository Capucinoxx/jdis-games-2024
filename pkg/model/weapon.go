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
  p.restrictToSquare(pos)
  p.setup(pos, config.ProjectileSize)

  return p
} 


// ApplyMovement updates the projectile's position based on its direction and a delta time.
func (p *Projectile) moveToDestination(dt float64) bool {
  dest := p.Destination

  dx := float64(dest.X - p.Position.X)
  dy := float64(dest.Y - p.Position.Y)
  dist := math.Abs(dx) + math.Abs(dy)

  if dist > config.ProjectileSpeed * float64(dt) {
    nextX := p.Position.X + dx/dist * config.ProjectileSpeed * dt
    nextY := p.Position.Y + dy/dist * config.ProjectileSpeed * dt

    p.Position.X = nextX
    p.Position.Y = nextY
    p.collider.ChangePosition(nextX, nextY)

    return true
  } else {
    p.Remove()
    return false
  } 
}


func (p *Projectile) restrictToSquare(pos *Point) {
	squareSize := 10.0

	left := math.Floor(pos.X/squareSize)*squareSize + (config.ProjectileSize / 2)
	right := left + squareSize - config.ProjectileSize
	top := math.Floor(pos.Y/squareSize)*squareSize + (config.ProjectileSize / 2)
	bottom := top + squareSize - config.ProjectileSize

	newX := p.Destination.X
	newY := p.Destination.Y

	originalDestX := p.Destination.X
	originalDestY := p.Destination.Y

	if originalDestX < left {
		newX = left
		newY = pos.Y + (left-pos.X)*(originalDestY-pos.Y)/(originalDestX-pos.X)
	} else if originalDestX > right {
		newX = right
		newY = pos.Y + (right-pos.X)*(originalDestY-pos.Y)/(originalDestX-pos.X)
	}

	if newY < top {
		newY = top
		newX = pos.X + (top-pos.Y)*(originalDestX-pos.X)/(originalDestY-pos.Y)
	} else if newY > bottom {
		newY = bottom
		newX = pos.X + (bottom-pos.Y)*(originalDestX-pos.X)/(originalDestY-pos.Y)
	}

	p.Destination.X = newX
	p.Destination.Y = newY
}


// IsCollidingWithEnvironment checks if the projectile is colliding with any non-projectile colliders in the map.
func (p *Projectile) IsCollidingWithEnvironment(m Map) bool {
  for _, collider := range m.Colliders() {	
    if collider.Type == ColliderProjectile {
      continue
		}
   
    if PolygonsIntersect(p.collider.polygon(), collider.polygon()) {
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
func (c *Cannon) Update(players []*Player, m Map, dt float64) {
	for _, p := range c.Projectiles {
		if !p.moveToDestination(dt) {
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
    }

    if p.IsCollidingWithEnvironment(m) {
			p.Remove()
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
  collider := c.owner.Collider()
  c.Projectiles = append(c.Projectiles, NewProjectile(
    &Point{X: collider.Pivot.X, Y: collider.Pivot.Y},
    &Point{X: pos.X, Y: pos.Y},
  ))
}

type Blade struct {
  collider *RectCollider
  owner *Player
}

func NewBlade(owner *Player) *Blade {
  blade := &Blade{ owner: owner }
  pivot := owner.Collider().Pivot
  blade.collider = NewRectCollider(pivot.X, pivot.Y, config.BladeSize)
  blade.collider.SetPivot(pivot.X, pivot.Y)
  return blade
}

func (b *Blade) Update(players []*Player, m Map, dt float64) {
  pivot := b.owner.Collider().Pivot
  b.collider.SetPivot(pivot.X, pivot.Y)
  b.collider.Rotate(config.BladeRotationSpeed * dt)


  for _, enemy := range players {
    if b.owner.Nickname == enemy.Nickname {
      continue
    }

    if PolygonsIntersect(b.collider.polygon(), enemy.Collider().polygon()) {
      enemy.TakeDmg(config.BladeDmg)
    }
  }
}

