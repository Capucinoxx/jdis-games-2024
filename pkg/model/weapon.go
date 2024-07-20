package model

import (
	"math"

	"github.com/capucinoxx/forlorn/consts"
)

// Projectile represents a moving projectile in the game.
type Projectile struct {
	Object
	ttl         float64
	Destination *Point
}

func NewProjectile(pos *Point, dest *Point) *Projectile {
	p := &Projectile{Destination: dest, ttl: consts.ProjectileTTL}
	p.setup(pos, consts.ProjectileSize)

	return p
}

// reduceTTL decrements the projectile's time to live by a delta time.
func (p *Projectile) reduceTTL(dt float64) {
	p.ttl -= dt
	if p.ttl <= 0 {
		p.Remove()
	}
}

// moveToDestination updates the projectile's position based on its direction and a delta time.
func (p *Projectile) moveToDestination(dt float64) {
	dest := p.Destination

	dx := float64(dest.X - p.Position.X)
	dy := float64(dest.Y - p.Position.Y)
	dist := math.Abs(dx) + math.Abs(dy)

	if dist > consts.ProjectileSpeed*float64(dt) {
		nextX := p.Position.X + dx/dist*consts.ProjectileSpeed*dt
		nextY := p.Position.Y + dy/dist*consts.ProjectileSpeed*dt

		p.Position.X = nextX
		p.Position.Y = nextY
		p.collider.ChangePosition(nextX, nextY)
	} else {
		p.Remove()
	}
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
		p.reduceTTL(dt)
		p.moveToDestination(dt)

		for _, enemy := range players {
			if c.owner.Nickname == enemy.Nickname {
				continue
			}

			if p.IsCollidingWithPlayer(enemy) {
				enemy.TakeDmg(consts.ProjectileDmg)
				p.Remove()
				continue
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
	collider := c.owner.Collider()
	c.Projectiles = append(c.Projectiles, NewProjectile(
		&Point{X: collider.Pivot.X, Y: collider.Pivot.Y},
		&Point{X: pos.X, Y: pos.Y},
	))
}

type Blade struct {
	collider *RectCollider
	owner    *Player
}

func NewBlade(owner *Player) *Blade {
	blade := &Blade{owner: owner}
	pivot := owner.Collider().Pivot
	blade.collider = NewRectLineCollider(pivot.X, pivot.Y, consts.BladeSize, consts.PlayerSize/4.0)
	blade.collider.SetPivot(pivot.X, pivot.Y)
	return blade
}

func (b *Blade) Update(players []*Player, m Map, rotation *float64) {
	pivot := b.owner.Collider().Pivot
	b.collider.SetPivot(pivot.X, pivot.Y)

	if rotation != nil {
		b.collider.Rotate(*rotation)
	}

	for _, enemy := range players {
		if b.owner.Nickname == enemy.Nickname {
			continue
		}

		if PolygonsIntersect(b.collider.polygon(), enemy.Collider().polygon()) {
			enemy.TakeDmg(consts.BladeDmg)
		}
	}
}
