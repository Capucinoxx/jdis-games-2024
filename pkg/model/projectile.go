package model

const PROJECTILE_DMG = 30

// Projectile represents a moving projectile in the game.
type Projectile struct {
	Position  *Point
	Direction *Point

	// // cleanup indicates whether the projectile should be removed from the game.
	cleanup bool
}

// ApplyMovement updates the projectile's position based on its direction and a delta time.
func (p *Projectile) ApplyMovement(dt float32) {
	p.Position.X += p.Direction.X * 100 * dt
	p.Position.Y += p.Direction.Y * 100 * dt
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
func (p *Projectile) IsAlive() bool { return p.cleanup }

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
		p.ApplyMovement(dt)

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
	projectile := &Projectile{
		Position: &Point{X: pos.X, Y: pos.Y},
	}
	projectile.Direction = projectile.Position.DirectionTo(c.owner.Collider.Pivot)
	c.Projectiles = append(c.Projectiles, projectile)
}
