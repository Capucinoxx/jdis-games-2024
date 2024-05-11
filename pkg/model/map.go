package model

import (
	"fmt"
	"math"
	"math/rand"
)

// Point représente un point "continu" dans un espace 2D.
type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

func NullPoint() *Point {
	return &Point{X: 0, Y: 0}
}

// Directions représente les directions possibles dans un espace 2D.
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

// DirectionTo retourne un vecteur normalisé pointant vers la destination à
// partir du point actuel.
func (p *Point) DirectionTo(dest *Point) *Point {
	dir := &Point{X: 0, Y: 0}

	dir.X = p.X - dest.X
	dir.Y = p.Y - dest.Y
	dir.normalize()

	return dir
}

// Add ajoute le vecteur spécifié à ce vecteur et retourne le résultat.
func (p *Point) Add(other *Point) *Point {
	return &Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Reflect retourne le vecteur réfléchi par rapport à la normale spécifiée.
func (p *Point) Reflect(normal *Point) *Point {
	dot := 2 * (p.X*normal.X + p.Y*normal.Y)
	return &Point{
		X: p.X - dot*normal.X,
		Y: p.Y - dot*normal.Y,
	}
}

// WithinDistanceOf retourne vrai si le point est à une distance inférieure au
// rayon spécifié de l'autre point. Sinon, il retourne faux.
func (p *Point) WithinDistanceOf(radius float32, oth *Point) bool {
	return (math.Pow(float64(oth.X-p.X), 2) +
		math.Pow(float64(oth.Y-p.Y), 2)) < math.Pow(float64(radius), 2)
}

// IsInPolygon retourne vrai si le point est à l'intérieur du polygone spécifié.
// Sinon, il retourne faux.
func (p *Point) IsInPolygon(poly []*Point) bool {
	inside := false
	for i, j := 0, len(poly)-1; i < len(poly); i++ {
		pi, pj := poly[i], poly[j]

		// regarde si est entre les ordonnées de pi et pj
		cond := (pi.Y > p.Y) != (pj.Y > p.Y)

		// interpoler l'abscisse de l'intersection entre la droite horizontale
		// passant par p et la droite passant par pi et pj à l'ordonnée p.Y
		px := (pj.X-pi.X)*(p.Y-pi.Y)/(pj.Y-pi.Y) + pi.X

		if cond && p.X < px {
			inside = !inside
		}
		j = i
	}
	return inside
}

func (p *Point) Hash() uint64 {
	return uint64(math.Float32bits(p.X))<<32 | uint64(math.Float32bits(p.Y))
}

// normalize normalise le vecteur. Cela signifie que la longueur du vecteur est
// égale à 1.
func (p *Point) normalize() {
	length := math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
	p.X /= float32(length)
	p.Y /= float32(length)
}

type ColliderType uint8

const (
	ColliderWall ColliderType = iota
	ColliderProjectile
)

// Collider représente un polygone qui représentant une collision dans le jeu.
type Collider struct {
	Points []*Point     `json:"points"`
	Type   ColliderType `json:"type"`
}

// polygon retourne le polygone représenté par le Collider.
func (c *Collider) polygon() Polygon {
	return Polygon{points: c.Points}
}

type Grid struct {
	height, width int
	cells         map[Point]map[Point]bool
}

// IsInBounds retourne vrai si le point est à l'intérieur des limites de la carte.
// Sinon, retourne faux.
func (g *Grid) isInBounds(pos *Point) bool {
	return pos.X >= 0 && pos.X < float32(g.width) && pos.Y >= 0 && pos.Y < float32(g.height)
}

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

// Map représente une carte de jeu. Elle contient des informations sur les
// collisions et les points de spawn.
type Map struct {
	Colliders []*Collider
	Spawns    []*Point
	cellSize  float32
}

// Populate remplit la carte avec des collisions en utilisant la grille spécifiée.
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

// Clear supprime toutes les collisions et les points de spawn de la carte.
func (m *Map) Clear() {
	m.Colliders = nil
	m.Spawns = nil
}
