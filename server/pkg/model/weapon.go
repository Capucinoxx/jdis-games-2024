package model

import (
	"math"

	"github.com/capucinoxx/jdis-games-2024/consts"
	"github.com/capucinoxx/jdis-games-2024/pkg/utils"
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

	dx := dest.X - p.Position.X
	dy := dest.Y - p.Position.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > consts.ProjectileSpeed*dt {
		nextX := p.Position.X + dx/dist*consts.ProjectileSpeed*dt
		nextY := p.Position.Y + dy/dist*consts.ProjectileSpeed*dt

		p.Position.X = nextX
		p.Position.Y = nextY
		p.collider.ChangePosition(nextX, nextY)
	} else {
		p.Position.X = dest.X
		p.Position.Y = dest.Y
		p.collider.ChangePosition(dest.X, dest.Y)
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
func (c *Cannon) Update(players []*Player, dt float64) {
	for _, p := range c.Projectiles {
		p.reduceTTL(dt)
		p.moveToDestination(dt)

		for _, enemy := range players {
			if c.owner.Nickname == enemy.Nickname || !enemy.IsAlive() {
				continue
			}

			if p.IsCollidingWithPlayer(enemy) {
				enemy.TakeDmg(consts.ProjectileDmg)
				c.owner.score += consts.ScoreOnHitWithProjectile
				p.Remove()

				utils.Log(c.owner.Nickname, "score", "hit %s with projectile +%d total: %d",
					enemy.Nickname, consts.ScoreOnHitWithProjectile, c.owner.score)
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

func (b *Blade) Update(players []*Player, rotation *float64) {
	pivot := b.owner.Collider().Pivot
	b.collider.ChangePosition(pivot.X, pivot.Y)

	if rotation != nil {
		b.collider.Rotate(*rotation)
	}

	for _, enemy := range players {
		if b.owner.Nickname == enemy.Nickname || !enemy.IsAlive() {
			continue
		}

		if PolygonsIntersect(b.collider.polygon(), enemy.Collider().polygon()) {
			enemy.TakeDmg(consts.BladeDmg)
			b.owner.score += consts.ScoreOnHitWithBlade

			utils.Log(b.owner.Nickname, "score", "hit %s with blade +%d total: %d",
				enemy.Nickname, consts.ScoreOnHitWithBlade, b.owner.score)
		}
	}
}
