package model

import (
	"math/rand"

	"github.com/capucinoxx/forlorn/pkg/config"
)

type Scorer struct {
  Object
  Value int
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


