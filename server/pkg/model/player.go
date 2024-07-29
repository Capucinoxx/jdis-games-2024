package model

import (
	"encoding/base64"
	"math"
	"sync"
	"time"

	"github.com/capucinoxx/jdis-games-2024/consts"
	"github.com/capucinoxx/jdis-games-2024/pkg/codec"
	"github.com/capucinoxx/jdis-games-2024/pkg/utils"
)

// Connection represents a network connection. It can be used for reading and writing data over the network
// and also for identifying the connection. The connection must be closed after use.
type Connection interface {
	// Identifier returns a unique identifier for the connection.
	Identifier() string

	// Close terminates the connection after a specified timeout.
	Close(time.Duration, bool)

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

	IsAdmin() bool

	SetAdmin(bool)
}

type PlayerWeapon = int

// Controls struct represents the player's controls.
// When a control is activated, the player performs the corresponding action.
type Controls struct {
	Dest *Point  `json:"dest,omitempty"`
	Save *string `json:"save,omitempty"`

	SwitchWeapon *PlayerWeapon `json:"switch,omitempty"`
	Shoot        *Point        `json:"shoot,omitempty"`
	RotateBlade  *float64      `json:"rotate_blade,omitempty"`
}

const (
	PlayerWeaponNone (PlayerWeapon) = iota
	PlayerWeaponCanon
	PlayerWeaponBlade
)

type PlayerScore struct {
	Name  string
	Score int
}

type Player struct {
	Object
	Destination *Point

	Nickname         string
	Color            int
	Client           *Client
	health           int
	respawnCountdown float64

	Controls Controls

	currentWeapon PlayerWeapon
	cannon        *Cannon
	blade         *Blade
	score         int

	storage [100]byte
	mu      sync.RWMutex
}

func NewPlayer(name string, color int, pos *Point, conn Connection) *Player {
	p := &Player{
		Nickname: name,
		Color:    color,
		Client: &Client{
			Out:        make(chan []byte, 10),
			In:         make(chan ClientMessage, 10),
			connection: conn,
		},
		currentWeapon: PlayerWeaponNone,

		health: 100,
	}

	p.setup(pos, consts.PlayerSize)
	p.cannon = NewCanon(p)
	p.blade = NewBlade(p)

	return p
}

func (p *Player) Collider() *RectCollider {
	return p.collider
}

func (p *Player) TakeDmg(dmg int) {
	alive := p.IsAlive()
	p.health -= dmg

	if p.health < 0 && alive {
		p.Client.SetBlind(true)
	}
}

func (p *Player) AddScore(score int) {
	p.score += score
}

func (p *Player) Score() int {
	return p.score
}

func (p *Player) Storage() [100]byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.storage
}

func (p *Player) ClearStorage() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.storage = [100]byte{}
}

func (p *Player) IsAlive() bool {
	return p.health > 0
}

func (p *Player) Update(players []*Player, game *GameState, dt float64) {
	if !p.IsAlive() {
		p.respawnCountdown += dt
		return
	}

	p.HandleMovement(players, game.Map, dt)
	p.HandleWeapon(players, game.Map, dt)
	p.HandleCoinCollision(game.coins.List())
	p.HandleSave()
}

func (p *Player) HandleSave() {
	if p.Controls.Save == nil {
		return
	}

	bytes, err := base64.StdEncoding.DecodeString(*p.Controls.Save)
	if err != nil {
		return
	}

	p.Controls.Save = nil

	p.mu.Lock()
	copy(p.storage[:], bytes)
	p.mu.Unlock()
}

func (p *Player) HandleCoinCollision(coins []*Scorer) {
	for _, coin := range coins {
		if coin.IsCollidingWithPlayer(p) {
			utils.Log(p.Nickname, "score", "take coin +%d total: %d",
				coin.Value, p.score)
		}
	}
}

func (p *Player) HandleRespawn(game *GameState) {
	if !p.IsAlive() && p.respawnCountdown > consts.RespawnTime {
		p.Respawn(game)
	}
}

func (p *Player) Respawn(game *GameState) {
	p.health = 100
	p.respawnCountdown = 0
	p.Position = game.GetSpawnPoint()
	p.collider.ChangePosition(p.Position.X, p.Position.Y)

	p.blade.collider.Rotation = 0.0
	p.Client.SetBlind(false)
}

func (p *Player) HandleMovement(players []*Player, m Map, dt float64) {
	if p.Controls.Dest == nil {
		return
	}

	px, py := p.Position.X, p.Position.Y

	p.moveToDestination(dt)

	for _, collider := range m.Colliders() {
		if PolygonsIntersect(p.collider.polygon(), collider.polygon()) {
			p.Position.X, p.Position.Y = px, py
			p.collider.ChangePosition(px, py)
			return
		}
	}
}

func (p *Player) moveToDestination(dt float64) {
	dest := p.Controls.Dest

	dx := float64(dest.X - p.Position.X)
	dy := float64(dest.Y - p.Position.Y)
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > consts.PlayerSpeed*float64(dt) {
		nextX := p.Position.X + dx/dist*consts.PlayerSpeed*dt
		nextY := p.Position.Y + dy/dist*consts.PlayerSpeed*dt

		p.Position.X = nextX
		p.Position.Y = nextY
		p.collider.ChangePosition(nextX, nextY)
	}
}

func (p *Player) HandleWeapon(players []*Player, m Map, dt float64) {
	p.cannon.Update(players, dt)
	bladeCondition := p.Controls.SwitchWeapon == nil && p.currentWeapon == PlayerWeaponBlade
	p.blade.Update(players, utils.NilIf(p.Controls.RotateBlade, !bladeCondition))

	if p.Controls.SwitchWeapon != nil {
		p.currentWeapon = *p.Controls.SwitchWeapon
		return
	}

	if p.currentWeapon == PlayerWeaponCanon && p.Controls.Shoot != nil {
		p.cannon.ShootAt(*p.Controls.Shoot)
		p.Controls.Shoot = nil
	}

	p.Controls.RotateBlade = nil
}

type PlayerInfo struct {
	Nickname      string
	Color         int32
	Health        int32
	Score         int64
	Pos           Point
	Dest          *Point
	CurrentWeapon PlayerWeapon
	Projectiles   []struct {
		Uuid [16]byte
		Pos  Point
		Dest Point
	}
	Blade struct {
		Start    Point
		End      Point
		Rotation float64
	}
}

func (p *Player) Encode(w codec.Writer) (err error) {
	if err = w.WriteString(p.Nickname); err != nil {
		return
	}

	if err = w.WriteInt32(int32(p.Color)); err != nil {
		return
	}

	if err = w.WriteInt32(int32(p.health)); err != nil {
		return
	}

	if err = w.WriteInt64(int64(p.score)); err != nil {
		return
	}

	if err = p.Position.Encode(w); err != nil {
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

	if err = w.WriteUint8(uint8(p.currentWeapon)); err != nil {
		return
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

	// encode blade
	if err = p.blade.collider.rect.a.Encode(w); err != nil {
		return
	}

	if err = p.blade.collider.rect.b.Encode(w); err != nil {
		return
	}

	if err = w.WriteFloat64(p.blade.collider.Rotation); err != nil {
		return
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

	if p.Score, err = r.ReadInt64(); err != nil {
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

	var currentWeapon uint8
	if currentWeapon, err = r.ReadUint8(); err != nil {
		return
	}
	p.CurrentWeapon = PlayerWeapon(currentWeapon)

	var length int32
	if length, err = r.ReadInt32(); err != nil {
		return
	}

	p.Projectiles = make([]struct {
		Uuid [16]byte
		Pos  Point
		Dest Point
	}, length)
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

	// decode Blade
	p.Blade.Start = Point{}
	if err = p.Blade.Start.Decode(r); err != nil {
		return
	}

	p.Blade.End = Point{}
	if err = p.Blade.End.Decode(r); err != nil {
		return
	}

	if p.Blade.Rotation, err = r.ReadFloat64(); err != nil {
		return
	}

	return
}

type Client struct {
	Out        chan []byte
	In         chan ClientMessage
	connection Connection
	blind      bool
	mu         sync.RWMutex
}

func (c *Client) GetConnection() Connection {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connection
}

func (c *Client) Disconnect() {
	utils.SafeClose(c.Out)
	utils.SafeClose(c.In)
}

func (c *Client) SetConnection(conn Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connection = conn
}

func (c *Client) IsBlind() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.blind
}

func (c *Client) SetBlind(b bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.blind = b
}

type ClientMessage struct {
	MessageType MessageType
	Body        interface{}
}
