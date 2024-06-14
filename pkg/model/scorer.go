package model

import "github.com/capucinoxx/forlorn/pkg/config"

type Scorer struct {
  Object
  Value int
}

func NewCoin(pos *Point) *Scorer {
  s := &Scorer{Value: config.CoinValue}
  s.setup(pos, config.CoinSize)

  return s
}


