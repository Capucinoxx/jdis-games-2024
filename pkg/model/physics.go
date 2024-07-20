package model

import (
	"math"

	"github.com/capucinoxx/forlorn/pkg/codec"
)

// Point represents a continuous point in 2D space.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p *Point) Encode(w codec.Writer) (err error) {
	if err = w.WriteFloat64(p.X); err != nil {
		return
	}

	if err = w.WriteFloat64(p.Y); err != nil {
		return
	}

	return
}

func (p *Point) Decode(r codec.Reader) (err error) {
	if p.X, err = r.ReadFloat64(); err != nil {
		return
	}

	if p.Y, err = r.ReadFloat64(); err != nil {
		return
	}

	return
}

var Directions = []Point{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

var (
	North = Directions[0]
	South = Directions[1]
	Est   = Directions[2]
	West  = Directions[3]
)

var OppositeDirections = map[Point]Point{
	North: South,
	South: North,
	Est:   West,
	West:  Est,
}

func DirectionToIndex(p Point) int {
	for i, d := range Directions {
		if d.X == p.X && d.Y == p.Y {
			return i
		}
	}

	return -1
}

// DirectionTo returns a normalized vector pointing towards the destination from the
// current point.
func (p *Point) DirectionTo(dest *Point) *Point {
	dir := &Point{X: dest.X - p.X, Y: dest.Y - p.Y}
	dir.normalize()
	return dir
}

func (p *Point) Equals(oth *Point, tolerance float32) bool {
	return math.Abs(float64(p.X-oth.X)) <= float64(tolerance) &&
		math.Abs(float64(p.Y-oth.Y)) <= float64(tolerance)
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
	p.X /= length
	p.Y /= length
}

func Normalize(p Point) Point {
	length := math.Sqrt(float64(p.X*p.X) + float64(p.Y*p.Y))
	return Point{X: p.X / length, Y: p.Y / length}
}

type Rect struct {
	a, b, c, d *Point
}

type Polygon struct {
	vertices []*Point
}

type RectCollider struct {
	rect  *Rect
	Pivot *Point

	Rotation float64
}

func NewRectCollider(x, y, size float64) *RectCollider {
	return &RectCollider{
		rect: &Rect{
			a: &Point{X: x - size/2, Y: y + size/2},
			b: &Point{X: x + size/2, Y: y + size/2},
			c: &Point{X: x + size/2, Y: y - size/2},
			d: &Point{X: x - size/2, Y: y - size/2},
		},

		Pivot:    &Point{X: x, Y: y},
		Rotation: 0,
	}
}

func NewRectLineCollider(x, y, width, height float64) *RectCollider {
	return &RectCollider{
		rect: &Rect{
			a: &Point{X: x, Y: y + height/2},
			b: &Point{X: x + width, Y: y + height/2},
			c: &Point{X: x + width, Y: y - height/2},
			d: &Point{X: x, Y: y - height/2},
		},

		Pivot:    &Point{X: x, Y: y},
		Rotation: 0,
	}
}

func (r *RectCollider) SetPivot(x, y float64) {
	r.Pivot.X = x
	r.Pivot.Y = y
}

func (r *RectCollider) ChangePosition(px, py float64) {
	dx := px - r.Pivot.X
	dy := py - r.Pivot.Y

	r.rect.a.X += dx
	r.rect.a.Y += dy

	r.rect.b.X += dx
	r.rect.b.Y += dy

	r.rect.c.X += dx
	r.rect.c.Y += dy

	r.rect.d.X += dx
	r.rect.d.Y += dy

	r.Pivot.X = px
	r.Pivot.Y = py
}

func (r *RectCollider) polygon() Polygon {
	return Polygon{
		vertices: []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d},
	}
}

func (r *RectCollider) Collisions(oth Polygon) bool {
	return PolygonsIntersect(r.polygon(), oth)
}

func ProjectPolygon(axis Point, polygon Polygon) (float64, float64) {
	min := (polygon.vertices[0].X * axis.X) + (polygon.vertices[0].Y * axis.Y)
	max := min

	for _, vertex := range polygon.vertices {
		projection := (vertex.X * axis.X) + (vertex.Y * axis.Y)
		if projection < min {
			min = projection
		}
		if projection > max {
			max = projection
		}
	}

	return min, max
}

func (r *RectCollider) Rotate(theta float64) {
	r.rotate(theta, r.rect.a)
	r.rotate(theta, r.rect.b)
	r.rotate(theta, r.rect.c)
	r.rotate(theta, r.rect.d)

	r.Rotation = math.Mod((r.Rotation + theta), 360.0)
}

func (r *RectCollider) rotate(theta float64, p *Point) {
	sin := math.Sin(theta)
	cos := math.Cos(theta)

	p.X -= r.Pivot.X
	p.Y -= r.Pivot.Y

	x := p.X*cos - p.Y*sin
	y := p.X*sin + p.Y*cos

	p.X = x + r.Pivot.X
	p.Y = y + r.Pivot.Y
}

func PolygonsIntersect(a, b Polygon) bool {
	for i := 0; i < len(a.vertices); i++ {
		j := (i + 1) % len(a.vertices)
		edge := Point{
			X: a.vertices[j].X - a.vertices[i].X,
			Y: a.vertices[j].Y - a.vertices[i].Y,
		}
		axis := Normalize(Point{X: -edge.Y, Y: edge.X})

		min1, max1 := ProjectPolygon(axis, a)
		min2, max2 := ProjectPolygon(axis, b)

		if max1 < min2 || max2 < min1 {
			return false
		}
	}

	for i := 0; i < len(b.vertices); i++ {
		j := (i + 1) % len(b.vertices)
		edge := Point{
			X: b.vertices[j].X - b.vertices[i].X,
			Y: b.vertices[j].Y - b.vertices[i].Y,
		}
		axis := Normalize(Point{X: -edge.Y, Y: edge.X})

		min1, max1 := ProjectPolygon(axis, a)
		min2, max2 := ProjectPolygon(axis, b)

		if max1 < min2 || max2 < min1 {
			return false
		}
	}

	return true
}
