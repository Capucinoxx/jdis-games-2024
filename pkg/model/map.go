package model

import (
	"fmt"
	"math"
	"math/rand"
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

// polygon returns the polygon represented by the Collider.
func (c *Collider) polygon() Polygon {
	return Polygon{points: c.Points}
}

// Grid represents a 2D grid structure in a game environment.
type Grid struct {
	height, width int
	cells         map[Point]map[Point]bool
}

// isInBounds returns true if the point is within the grid's boundaries, otherwise false.
func (g *Grid) isInBounds(pos *Point) bool {
	return pos.X >= 0 && pos.X < float32(g.width) && pos.Y >= 0 && pos.Y < float32(g.height)
}

// GenerateGrid creates and returns a new Grid of the specified dimensions.
func GenerateGrid(width, height int) *Grid {
	grid := &Grid{
		height: height,
		width:  width,
		cells:  make(map[Point]map[Point]bool),
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grid.cells[Point{X: float32(x), Y: float32(y)}] = make(map[Point]bool)
		}
	}

	visited := make(map[Point]bool)

	var dfs func(Point)
	dfs = func(pos Point) {
		visited[pos] = true

		dirs := make([]Point, len(Directions))
		copy(dirs, Directions)
		rand.Shuffle(len(dirs), func(i, j int) {
			dirs[i], dirs[j] = dirs[j], dirs[i]
		})

		for _, dir := range dirs {
			v := dir.Add(&pos)
			if grid.isInBounds(v) && !visited[*v] {
				dfs(*v)
				grid.cells[pos][dir] = true
				grid.cells[*v][*dir.Reflect(NullPoint())] = true
			}
		}
	}

	dfs(Point{X: 0, Y: 0})
	random_pos := Point{X: float32(rand.Intn(width)), Y: float32(rand.Intn(height))}
	dfs(random_pos)

	return grid
}

// Map represents a game map, containing information about collisions and spawn points.
type Map struct {
	Colliders []*Collider
	Spawns    []*Point
	cellSize  float32
}

// Populate fills the map with colliders using the specified grid.
func (m *Map) Populate(grid *Grid) {
	// TODO: Fix It
	if m.cellSize == 0 {
		m.cellSize = 1
	}

	m.Colliders = append(m.Colliders, &Collider{
		Points: []*Point{{X: 0, Y: 0}, {X: 0, Y: float32(grid.height) * m.cellSize}},
	})

	m.Colliders = append(m.Colliders, &Collider{
		Points: []*Point{{X: 0, Y: 0}, {X: float32(grid.height) * m.cellSize, Y: 0}},
	})

	m.Colliders = append(m.Colliders, &Collider{
		Points: []*Point{{X: float32(grid.width) * m.cellSize, Y: 0}, {X: float32(grid.width) * m.cellSize, Y: float32(grid.height) * m.cellSize}},
	})

	m.Colliders = append(m.Colliders, &Collider{
		Points: []*Point{{X: 0, Y: float32(grid.height) * m.cellSize}, {X: float32(grid.width) * m.cellSize, Y: float32(grid.height) * m.cellSize}},
	})

	for y := 0; y < grid.height; y++ {
		for x := 0; x < grid.width; x++ {
			cell := grid.cells[Point{X: float32(x), Y: float32(y)}]
			x1, y1 := float32(x)*m.cellSize, float32(y)*m.cellSize
			x2, y2 := x1+m.cellSize, y1+m.cellSize

			if _, exist := cell[RIGHT]; exist {
				m.Colliders = append(m.Colliders, &Collider{
					Points: []*Point{{X: x2, Y: y1}, {X: x2, Y: y2}},
				})
			}

			if _, exist := cell[DOWN]; exist {
				m.Colliders = append(m.Colliders, &Collider{
					Points: []*Point{{X: x1, Y: y2}, {X: x2, Y: y2}},
				})
			}
		}
	}
}

// Clear removes all colliders and spawn points from the map.
func (m *Map) Clear() {
	m.Colliders = nil
	m.Spawns = nil
}
