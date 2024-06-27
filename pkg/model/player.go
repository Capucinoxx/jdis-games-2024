package model

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/config"
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


// Controls struct represents the player's controls.
// When a control is activated, the player performs the corresponding action.
type Controls struct {
  Dest *Point `json:"dest,omitempty"`

	// Shoot is the target point to shoot at, if any.
  Shoot *Point `json:"shoot,omitempty"`
}

// Player represents a player in the game.
type Player struct {
	Token            string
	Nickname         string
  Color int
	Health           atomic.Int32
	Score            float64
	respawnCountdown float64
	Client           *Client
	Controls         Controls
	Collider         *RectCollider
	cannon           *Cannon
}

// NewPlayer creates a new player with an initial position and a network connection.
func NewPlayer(name string, color int, x float64, y float64, conn Connection) *Player {
	p := &Player{
	  Nickname: name,
    Color: color,
		Collider: NewRectCollider(x, y, config.PlayerSize),
		Health:   atomic.Int32{},
		Client: &Client{
			Out:        make(chan []byte, 10),
			In:         make(chan ClientMessage, 10),
			Connection: conn,
		},
	}

	p.Health.Add(config.PlayerHealth)
	p.cannon = NewCanon(p)
	return p
}

// String returns a string representation of the player.
func (p *Player) String() string {
	return fmt.Sprintf("[%s: { pos: (%f, %f), v: %f, dest: %+v, health: %d }]", p.Nickname, p.Collider.Pivot.X, p.Collider.Pivot.Y, p.Collider.velocity, p.Controls, p.Health.Load())
}

// IsAlive returns true if the player's health is above zero, indicating they are alive.
func (p *Player) IsAlive() bool {
	return p.Health.Load() > 0
}

// Update updates the player's state based on the current game state.
func (p *Player) Update(players []*Player, game *GameState, dt float64) {
	m := game.Map
	if !p.IsAlive() {
		p.respawnCountdown += dt
		return
	}

	p.HandleMovement(players, m, dt)
	p.HandleCannon(players, m, dt)
}

// HandleMovement manages the player's movement based on their controls.
func (p *Player) HandleMovement(players []*Player, m Map, dt float64) {
  if p.Controls.Dest != nil {
    p.moveToDestination(players, m, dt)
  }
}

func (p *Player) moveToDestination(players []*Player, m Map, dt float64) {
  r := p.Collider
  dest := p.Controls.Dest

  dx := float64(dest.X - r.Pivot.X)
  dy := float64(dest.Y - r.Pivot.Y)
  dist := math.Abs(dx) + math.Abs(dy)
  
  speed := config.PlayerSpeed

  if dist > float64(speed*dt) {
    nextX := r.Pivot.X + dx/dist * speed * dt
    nextY := r.Pivot.Y + dy/dist * speed * dt

    if !p.checkCollisionAt(nextX, nextY, players, m) {
      r.Pivot.X = nextX
      r.Pivot.Y = nextY
    }
  } else {
    if !p.checkCollisionAt(dest.X, dest.Y, players, m) {
      r.Pivot.X = dest.X
      r.Pivot.Y = dest.Y
    }
    p.Controls.Dest = nil
  }
}

func (p *Player) checkCollisionAt(x, y float64, players []*Player, m Map) bool {
  originalX, originalY := p.Collider.Pivot.X, p.Collider.Pivot.Y
  p.Collider.Pivot.X, p.Collider.Pivot.Y = x, y

  collides := p.checkCollisionWithPlayers(players) || p.checkCollisionWithMap(m)

  p.Collider.Pivot.X, p.Collider.Pivot.Y = originalX, originalY
  return collides
}

// HandleCannon handles the player's cannon actions.
func (p *Player) HandleCannon(players []*Player, m Map, dt float64) {
	if p.Controls.Shoot != nil {
		p.cannon.ShootAt(*p.Controls.Shoot)
	}

  p.cannon.Update(players, m, dt)
}

func (p *Player) TakeDmg(dmg int32) {
	p.Health.Add(-dmg)
}

// checkCollisionWithPlayers checks if there is a collision with any other player.
func (p *Player) checkCollisionWithPlayers(players []*Player) bool {
	for _, ennemy := range players {
    if ennemy.Nickname == p.Nickname || !ennemy.IsAlive() {
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


func (p *Player) HandleRespawn(game *GameState) {
  if !p.IsAlive() && p.respawnCountdown > config.RespawnTime {
    spawn := game.GetSpawnPoint()
    p.Health.Store(config.PlayerHealth)
    p.respawnCountdown = 0

    p.Collider.ChangePosition(spawn.X, spawn.Y)
  }
}


// applyMovement applies the movement to the player based on their current direction and velocity.
func (p *Player) applyMovement() {
	r := p.Collider
	points := []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d, r.look, r.Pivot}
	for _, point := range points {
		point.X += r.Dir.X * r.velocity
		point.Y += r.Dir.Y * r.velocity
	}
}

type PlayerInfo struct {
  Nickname string
  Color int32
  Health int32
  Pos Point
  Dest *Point
  Projectiles []struct {
    Uuid [16]byte
    Pos Point
    Dest Point
  }
}

func (p *Player) Encode(w codec.Writer) (err error) {
  if err = w.WriteString(p.Nickname); err != nil {
    return
  }

  if err = w.WriteInt32(int32(p.Color)); err != nil {
    return
  }

  if err = w.WriteInt32(p.Health.Load()); err != nil {
    return
  }

  if err = p.Collider.Pivot.Encode(w); err != nil {
    return
  }

  if p.Controls.Dest != nil {
    if err = w.WriteBool(true); err != nil {
      return
    }

    if err = p.Controls.Dest.Encode(w); err != nil {
      return
    }
  } else {
    if err = w.WriteBool(false); err != nil {
      return
    }
  }

  bullets := p.cannon.Projectiles
  if err = w.WriteInt32(int32(len(bullets))); err != nil {
    return
  }

  for _, bullet := range bullets {
    if _, err = w.WriteBytes(bullet.uuid[:]); err != nil {
      return
    }

    if err = bullet.Position.Encode(w); err != nil {
      return
    }
    if err = bullet.Destination.Encode(w); err != nil {
      return
    }
  }

  return
}

func (p *PlayerInfo) Decode(r codec.Reader) (err error) {
  if p.Nickname, err = r.ReadString(); err != nil {
    return
  }

  if p.Color, err = r.ReadInt32(); err != nil {
    return
  }
  

  if p.Health, err = r.ReadInt32(); err != nil {
    return
  }

  
  if err = p.Pos.Decode(r); err != nil {
    return
  }

  var hasDest bool
  if hasDest, err = r.ReadBool(); err != nil {
    return
  }

  if hasDest {
    p.Dest = &Point{}
    if err = p.Dest.Decode(r); err != nil {
      return
    }
  }
  
  var length int32
  if length, err = r.ReadInt32(); err != nil {
    return
  }
  
  p.Projectiles = make([]struct{Uuid [16]byte; Pos Point; Dest Point}, length)
  for i := 0; i < int(length); i++ {
    var id []byte
    if id, err = r.ReadBytes(16); err != nil {
      return
    }
    copy(p.Projectiles[i].Uuid[:], id)

    if err = p.Projectiles[i].Pos.Decode(r); err != nil {
      return
    }

    if err = p.Projectiles[i].Dest.Decode(r); err != nil {
      return
    }
  }

  return
}

