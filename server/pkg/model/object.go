package model

import (
	"math/rand"

	"github.com/google/uuid"

	"github.com/capucinoxx/forlorn/consts"
	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/utils"
)

type Object struct {
	uuid     [16]byte
	Position *Point
	collider *RectCollider

	cleanup bool
}

func (o *Object) setup(pos *Point, size float64) {
	o.uuid = uuid.New()
	o.Position = pos
	o.collider = NewRectCollider(pos.X, pos.Y, size)
	o.cleanup = false
}

func (o *Object) IsCollidingWithPlayer(player *Player) bool {
	return o.collider.Collisions(player.Collider().polygon())
}

func (o *Object) Remove() { o.cleanup = true }

func (o *Object) IsAlive() bool { return !o.cleanup }

type Scorer struct {
	Object
	Value int32
}

func NewCoin(pos ...*Point) *Scorer {
	if len(pos) == 0 {
		pos = []*Point{{
			X: rand.Float64() * float64(consts.MapWidth*consts.CellWidth),
			Y: rand.Float64() * float64(consts.MapWidth*consts.CellWidth),
		}}
	}
	s := &Scorer{Value: consts.CoinValue}

	s.setup(pos[0], consts.CoinSize)

	return s
}

func NewBigCoin(center *Point) *Scorer {
	s := &Scorer{Value: consts.BigCoinValue}
	s.setup(center, consts.BigCoinSize)

	return s
}

func (s *Scorer) IsCollidingWithPlayer(player *Player) bool {
	if !s.IsAlive() {
		return false
	}

	ok := PolygonsIntersect(s.collider.polygon(), player.Collider().polygon())
	if ok {
		player.AddScore(int(s.Value))
		s.Remove()
	}

	return ok
}

func (s *Scorer) Encode(w codec.Writer) (err error) {
	if _, err = w.WriteBytes(s.uuid[:]); err != nil {
		return
	}

	if err = s.Position.Encode(w); err != nil {
		return
	}

	if err = w.WriteInt32(s.Value); err != nil {
		return
	}

	return
}

type Scorers struct {
	scorers []*Scorer
}

func NewScorers() *Scorers {
	return &Scorers{scorers: []*Scorer{}}
}

func (s *Scorers) Add(scorers ...*Scorer) {
	s.scorers = append(s.scorers, scorers...)
}

func (s *Scorers) Set(scorers []*Scorer) {
	s.scorers = scorers
}

func (s *Scorers) Update() bool {
	for i := 0; i < len(s.scorers); i++ {
		if !s.scorers[i].IsAlive() {
			if len(s.scorers) == 1 {
				return true
			}
			s.scorers[i] = NewCoin()
			utils.Log("coin", "score", "new coin spawn position (%f, %f)", s.scorers[i].Position.X, s.scorers[i].Position.Y)
		}
	}
	return false
}

func (s *Scorers) List() []*Scorer {
	return s.scorers
}
