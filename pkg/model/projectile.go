package model

type Projectile struct {
	Position  *Point
	Direction *Point
	cleanup   bool
}

func (p *Projectile) ApplyMovement(dt float32) {
	p.Position.X += p.Direction.X * 100 * dt
	p.Position.Y += p.Direction.Y * 100 * dt
}

func (p *Projectile) IsCollidingWithPlayer(player *Player) bool {
	rect := player.Collider.rect
	polygon := []*Point{rect.a, rect.b, rect.c, rect.d}
	return p.Position.IsInPolygon(polygon)
}

func (p *Projectile) IsCollidingWithEnvironment(m *Map) bool {
	for _, collider := range m.Colliders {
		if collider.Type == ColliderProjectile {
			continue
		}

		if p.Position.IsInPolygon(collider.Points) {
			return true
		}
	}

	return false
}

func (p *Projectile) Remove()       { p.cleanup = true }
func (p *Projectile) IsAlive() bool { return p.cleanup }

type Cannon struct {
	Projectiles []*Projectile
	owner       *Player
}

func NewCanon(owner *Player) *Cannon {
	return &Cannon{
		owner: owner,
	}
}

func (c *Cannon) Update(players []*Player, m *Map, dt float32) {
	for _, p := range c.Projectiles {
		p.ApplyMovement(dt)

		for _, enemy := range players {
			if c.owner.ID == enemy.ID {
				continue
			}

			if p.IsCollidingWithPlayer(enemy) {
				// TODO
			}

			if p.IsCollidingWithEnvironment(m) {
				// TODO
			}
		}
	}

	projectiles := make([]*Projectile, 0)
	for _, p := range c.Projectiles {
		if p.IsAlive() {
			projectiles = append(projectiles, p)
		}
	}
	c.Projectiles = projectiles
}

func (c *Cannon) ShootAt(pos Point) {

	projectile := &Projectile{
		Position: &Point{X: pos.X, Y: pos.Y},
	}
	projectile.Direction = projectile.Position.DirectionTo(c.owner.Collider.Pivot)
	c.Projectiles = append(c.Projectiles, projectile)
}
