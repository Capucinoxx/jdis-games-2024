package model

import (
	"github.com/capucinoxx/jdis-games-2024/pkg/codec"
)

// ColliderType defines the type of object in the game environment that can collide.
type ColliderType uint8

const (
	ColliderWall ColliderType = iota
	ColliderProjectile
)

// Collider represents a polygon that models a collision in the game.
type Collider struct {
	Points []*Point     `json:"points"`
	Type   ColliderType `json:"type"`
}

func (c *Collider) Encode(w codec.Writer) (err error) {
	if err = w.WriteUint8((uint8(len(c.Points)))); err != nil {
		return
	}

	for _, p := range c.Points {
		if err = p.Encode(w); err != nil {
			return
		}
	}

	if err = w.WriteUint8(uint8(c.Type)); err != nil {
		return
	}

	return
}

func (c *Collider) Decode(r codec.Reader) (err error) {
	var size uint8
	if size, err = r.ReadUint8(); err != nil {
		return
	}

	c.Points = make([]*Point, 0, size)
	for i := uint8(0); i < size; i++ {
		p := &Point{}
		if err = p.Decode(r); err != nil {
			return
		}
		c.Points = append(c.Points, p)
	}

	var cType uint8
	if cType, err = r.ReadUint8(); err != nil {
		return
	}

	c.Type = ColliderType(cType)
	return
}

// polygon returns the polygon represented by the Collider.
func (c *Collider) polygon() Polygon {
	return Polygon{vertices: c.Points}
}

// Map represents a game map, containing information about collisions and spawn points.
type Map interface {
	Setup()
	Centroid() Point
	Colliders() []*Collider
	Spawns(int) []*Point
	Size() int
	DiscreteMap() [][]uint8
	Encode(codec.Writer, bool) error
}
