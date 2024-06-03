package model

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/capucinoxx/forlorn/pkg/codec"
)

// Connection represents a network connection. It can be used for reading and writing data over the network
// and also for identifying the connection. The connection must be closed after use.
type Connection interface {
	// Identifier returns a unique identifier for the connection.
	Identifier() string

	// Close terminates the connection after a specified timeout.
	Close(time.Duration)

	// PrepareRead prepares the connection to read a specified amount of data within a timeout.
	PrepareRead(int64, time.Duration)

	// Read retrieves data from the connection.
	Read() ([]byte, error)

	// PrepareWrite prepares the connection for writing within a specified timeout.
	PrepareWrite(time.Duration)

	// Write sends data over the connection.
	Write([]byte) error

	// Ping sends a ping to check the connectivity and latency.
	Ping(time.Duration)
}

const (
	// playerSize defines the default size of a player.
	playerSize = 0.1

	// defaultHealth is the starting health of a player.
	defaultHealth = 100
)

// Controls struct represents the player's controls.
// When a control is activated, the player performs the corresponding action.
type Controls struct {
	Rotation uint32

	// Shoot is the target point to shoot at, if any.
	Shoot *Point
}

// Player represents a player in the game.
type Player struct {
	ID               uint8
	Token            string
	Nickname         string
	Health           atomic.Int32
	Score            float64
	respawnCountdown float32
	Client           *Client
	Controls         Controls
	Collider         *RectCollider
	cannon           *Cannon
}

// NewPlayer creates a new player with an initial position and a network connection.
func NewPlayer(id uint8, x float32, y float32, conn Connection) *Player {
	p := &Player{
		ID:       id,
		Collider: NewRectCollider(x, y, playerSize),
		Health:   atomic.Int32{},
		Client: &Client{
			Out:        make(chan []byte, 10),
			In:         make(chan ClientMessage, 10),
			Connection: conn,
		},
	}

	p.Health.Add(defaultHealth)
	p.cannon = NewCanon(p)
	return p
}

// String returns a string representation of the player.
func (p *Player) String() string {
	return fmt.Sprintf("[%d: { pos: (%f, %f), v: %f, rot: %d, health: %d }]", p.ID, p.Collider.Pivot.X, p.Collider.Pivot.Y, p.Collider.velocity, p.Collider.Rotation, p.Health.Load())
}

// IsAlive returns true if the player's health is above zero, indicating they are alive.
func (p *Player) IsAlive() bool {
	return p.Health.Load() > 0
}

// Update updates the player's state based on the current game state.
func (p *Player) Update(players []*Player, game *GameState, dt float32) {
	m := game.Map
	if !p.IsAlive() {
		p.respawnCountdown += dt
		return
	}

	p.HandleMovement(players, m, dt)
	p.HandleCannon(players, m, dt)
}

// HandleMovement manages the player's movement based on their controls.
func (p *Player) HandleMovement(players []*Player, m Map, dt float32) {
	r := p.Collider

	hasCollision := p.checkCollisionWithPlayers(players) || p.checkCollisionWithMap(m)

	p.updateVelocity(dt, hasCollision)
	p.updateRotation()
	if !hasCollision {
		p.applyMovement()
	}

	r.Rotation = (r.Rotation + p.Controls.Rotation) % 360
}

// HandleCannon handles the player's cannon actions.
func (p *Player) HandleCannon(players []*Player, m Map, dt float32) {
	if p.Controls.Shoot != nil {
		p.cannon.ShootAt(*p.Controls.Shoot)
	}
}

func (p *Player) TakeDmg(dmg int32) {
	p.Health.Add(-dmg)
}

// checkCollisionWithPlayers checks if there is a collision with any other player.
func (p *Player) checkCollisionWithPlayers(players []*Player) bool {
	for _, ennemy := range players {
		if ennemy.ID == p.ID || !ennemy.IsAlive() {
			continue
		}

		if p.Collider.Collisions(ennemy.Collider.polygon()) {
			return true
		}
	}

	return false
}

// checkCollisionWithMap checks if the player collides with the map.
func (p *Player) checkCollisionWithMap(m Map) bool {
	for _, collider := range m.Colliders() {
		if p.Collider.Collisions(collider.polygon()) {
			return true
		}
	}

	return false
}

// updateVelocity updates the player's velocity based on collision.
// If there is a collision, the velocity is set to zero.
func (p *Player) updateVelocity(dt float32, hasCollision bool) {
	r := p.Collider
	r.velocity = defaultForwardSpeed
}

// updateRotation updates the player's rotation based on their controls.
func (p *Player) updateRotation() {
	p.applyRotation(p.Controls.Rotation)
}

// applyRotation applies the specified rotation to the player.
func (p *Player) applyRotation(rd uint32) {
	r := p.Collider
	points := []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d, r.look}
	for _, point := range points {
		r.rotate(rd, point)
	}
	r.CalculDirection()
}

// applyMovement applies the movement to the player based on their current direction and velocity.
func (p *Player) applyMovement() {
	r := p.Collider
	points := []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d, r.look, r.Pivot}
	for _, point := range points {
		point.X += r.dir.X * r.velocity
		point.Y += r.dir.Y * r.velocity
	}
}

func (p *Player) Encode(w codec.Writer) (err error) {
  if err = w.WriteString(p.Nickname); err != nil {
    return
  }

  if err = w.WriteInt32(p.Health.Load()); err != nil {
    return
  }

  if err = p.Collider.Pivot.Encode(w); err != nil {
    return
  }

  bullets := p.cannon.Projectiles
  if err = w.WriteInt32(int32(len(bullets))); err != nil {
    return
  }

  for _, bullet := range bullets {
    if err = bullet.Position.Encode(w); err != nil {
      return
    }
    if err = bullet.Direction.Encode(w); err != nil {
      return
    }
  }

  return
}

func (p *Player) Decode(r codec.Reader) (err error) {
  if p.Nickname, err = r.ReadString(); err != nil {
    return
  }

  var health int32
  if health, err = r.ReadInt32(); err != nil {
    return
  }
  p.Health.Store(health)

  p.Collider = &RectCollider{Pivot: &Point{}}
  if err = p.Collider.Pivot.Decode(r); err != nil {
    return
  }
  
  var bullets_length int32
  if bullets_length, err = r.ReadInt32(); err != nil {
    return
  }

  p.cannon = &Cannon{Projectiles: make([]*Projectile, bullets_length)}
  for i := 0; i < int(bullets_length); i++ {
    pos := &Point{}
    if err = pos.Decode(r); err != nil {
      return
    }
    direction := &Point{}
    if err = pos.Decode(r); err != nil {
      return
    }

    p.cannon.Projectiles[i] = &Projectile{
      Position: pos,
      Direction: direction,
    }
  }

  return
}
