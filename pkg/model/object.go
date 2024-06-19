package model

import "github.com/google/uuid"

type Object struct {
  uuid [16]byte
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
  return o.collider.Collisions(player.Collider.polygon())
}

func (o *Object) Remove() { o.cleanup = true }

func (o *Object) IsAlive() bool { return !o.cleanup }
