package model

import (
	"math/rand"

	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/config"
)

type Scorer struct {
  Object
  Value int32
}

func NewCoin() *Scorer {
  s := &Scorer{Value: config.CoinValue}
  pos := &Point{
    X: rand.Float64() * float64(config.MapWidth),
    Y: rand.Float64() * float64(config.MapWidth),
  }

  s.setup(pos, config.CoinSize)

  return s
}

func (s *Scorer) IsCollidingWithPlayer(player *Player) bool {
  if !s.IsAlive() {
    return false
  }

  ok := PolygonsIntersect(s.collider.polygon(), player.Collider().polygon())
  if ok {
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

