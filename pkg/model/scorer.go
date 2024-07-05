package model

import (
  "math/rand"

  "github.com/capucinoxx/forlorn/consts"
  "github.com/capucinoxx/forlorn/pkg/codec"
)

type Scorer struct {
  Object
  Value int32
}

func NewCoin() *Scorer {
  s := &Scorer{Value: consts.CoinValue}
  pos := &Point{
    X: rand.Float64() * float64(consts.MapWidth * consts.CellWidth),
    Y: rand.Float64() * float64(consts.MapWidth * consts.CellWidth),
  }

  s.setup(pos, consts.CoinSize)

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
  return &Scorers{}
}

func (s *Scorers) Add(scorer *Scorer) {
  s.scorers = append(s.scorers, scorer)
}

func (s *Scorers) Update(players []*Player) {
  for _, scorer := range s.scorers {
    for _, player := range players {
      scorer.IsCollidingWithPlayer(player)
    }
  }

  for i := 0; i < len(s.scorers); i++ {
    if !s.scorers[i].IsAlive() {
      s.scorers[i] = NewCoin()
    }
  }
}
