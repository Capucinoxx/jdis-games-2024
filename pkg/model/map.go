package model

import (
	"fmt"
	"math"

	"github.com/capucinoxx/forlorn/pkg/codec"
)

// Point represents a continuous point in 2D space.
type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// String returns the string representation of the Point in the format "(X, Y)".
func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

func (p *Point) Encode(w codec.Writer) (err error) {
	if err = w.WriteFloat32(p.X); err != nil {
		return
	}

	if err = w.WriteFloat32(p.Y); err != nil {
		return
	}

	return
}

func (p *Point) Decode(r codec.Reader) (err error) {
	if p.X, err = r.ReadFloat32(); err != nil {
		return
	}

	if p.Y, err = r.ReadFloat32(); err != nil {
		return
	}

	return
}

// NullPoint returns a new Point with both coordinates set to zero.
func NullPoint() *Point {
	return &Point{X: 0, Y: 0}
}

// Directions represents the possible movement directions in a 2D space.
var Directions = []Point{
	{X: 0, Y: -1}, // Up
	{X: 1, Y: 0},  // Right
	{X: 0, Y: 1},  // Down
	{X: -1, Y: 0}, // Left
}

var (
	UP    = Directions[0]
	RIGHT = Directions[1]
	DOWN  = Directions[2]
	LEFT  = Directions[3]
)

// DirectionTo returns a normalized vector pointing towards the destination from the
// current point.
func (p *Point) DirectionTo(dest *Point) *Point {
	dir := &Point{X: dest.X - p.X, Y: dest.Y - p.Y}
	dir.normalize()
	return dir
}

// Add adds the specified vector to this vector and returns the result.
func (p *Point) Add(other *Point) *Point {
	return &Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Reflect returns the vector reflected by the specified normal.
func (p *Point) Reflect(normal *Point) *Point {
	dot := 2 * (p.X*normal.X + p.Y*normal.Y)
	return &Point{
		X: p.X - dot*normal.X,
		Y: p.Y - dot*normal.Y,
	}
}

// WithinDistanceOf returns true if the point is within the specified radius of another
// point, otherwise false.
func (p *Point) WithinDistanceOf(radius float32, oth *Point) bool {
	return (math.Pow(float64(oth.X-p.X), 2) +
		math.Pow(float64(oth.Y-p.Y), 2)) < math.Pow(float64(radius), 2)
}

// IsInPolygon returns true if the point is inside the specified polygon, otherwise false.
func (p *Point) IsInPolygon(poly []*Point) bool {
	inside := false
	for i, j := 0, len(poly)-1; i < len(poly); i++ {
		pi, pj := poly[i], poly[j]

		// Checks if it is between the y-coordinates of pi and pj.
		cond := (pi.Y > p.Y) != (pj.Y > p.Y)

		// Interpolate the x-coordinate of the intersection between the horizontal line through
		// p and the line passing through pi and pj at the y-coordinate p.Y.
		px := (pj.X-pi.X)*(p.Y-pi.Y)/(pj.Y-pi.Y) + pi.X

		if cond && p.X < px {
			inside = !inside
		}
		j = i
	}
	return inside
}

// normalize normalizes the vector such that its length becomes 1.
func (p *Point) normalize() {
	length := math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
	p.X /= float32(length)
	p.Y /= float32(length)
}

// ================================================================================================
// ================================================================================================

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
	if err = w.WriteInt32(int32(len(c.Points))); err != nil {
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
	var size int32
	if size, err = r.ReadInt32(); err != nil {
		return
	}

	for i := int32(0); i < size; i++ {
		var p Point
		if err = p.Decode(r); err != nil {
			return
		}
		c.Points = append(c.Points, &p)
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
	return Polygon{points: c.Points}
}

// Map represents a game map, containing information about collisions and spawn points.
type Map interface {
	Setup()
	Colliders() []*Collider
	Spawns() []*Point
}
