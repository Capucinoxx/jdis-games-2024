package model

import (
	"math"
)

const (
	// defaultForwardSpeed est la vitesse de déplacement par défaut.
	defaultForwardSpeed = 1
)

// Rect représente un rectangle dans un espace 2D.
type Rect struct {
	a, b, c, d *Point
}

// Polygon représente un polygone dans un espace 2D.
type Polygon struct {
	points []*Point
}

// String retourne une représentation en chaîne de caractères du polygone.
func (p Polygon) String() string {
	str := "["
	for _, point := range p.points {
		str += point.String() + ", "
	}

	str = str[:len(str)-2]
	str += "]"
	return str
}

// RectCollider représente un rectangle avec des capacités de collision
// dans un espace 2D.
type RectCollider struct {
	rect  *Rect
	look  *Point
	Dir   *Point
	Pivot *Point

	Rotation     uint32
	forwardSpeed float32

	velocity float32
}

// NewRectCollider crée un nouveau RectCollider.
func NewRectCollider(x, y, size float32) *RectCollider {
	return &RectCollider{
		rect: &Rect{
			a: &Point{X: x - size/2, Y: y + size/2},
			b: &Point{X: x + size/2, Y: y + size/2},
			c: &Point{X: x + size/2, Y: y - size/2},
			d: &Point{X: x - size/2, Y: y - size/2},
		},

		Pivot: &Point{X: x, Y: y},
		look:  &Point{X: x, Y: y + 2},
		Dir:   &Point{X: 0, Y: 0},

		Rotation:     0,
		forwardSpeed: defaultForwardSpeed,
	}
}

// CalculDirection calcule la direction du RectCollider et
// la normalise.
func (r *RectCollider) CalculDirection() {
	r.Dir.X = r.look.X - r.Pivot.X
	r.Dir.Y = r.look.Y - r.Pivot.Y
	r.Dir.normalize()
}

// Rotate tourne le RectCollider de l'angle spécifié.
func (r *RectCollider) Rotate(angle uint32) {
	r.rotate(angle, r.look)
	r.rotate(angle, r.rect.a)
	r.rotate(angle, r.rect.b)
	r.rotate(angle, r.rect.c)
	r.rotate(angle, r.rect.d)
}

// rotate tourne le point spécifié autour du pivot du RectCollider.
func (r *RectCollider) rotate(theta uint32, p *Point) {
	sint := float32(math.Sin(float64(theta)))
	cost := float32(math.Cos(float64(theta)))

	p.X -= r.Pivot.X
	p.Y -= r.Pivot.Y

	x := p.X*cost - p.Y*sint
	y := p.X*sint + p.Y*cost

	p.X = x + r.Pivot.X
	p.Y = y + r.Pivot.Y
}

// polygon retourne le polygone représenté par le RectCollider.
func (r *RectCollider) polygon() Polygon {
	return Polygon{
		points: []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d},
	}
}

// Collisions retourne vrai si le RectCollider entre en collision avec
// le polygone spécifié. Sinon retourne faux.
func (r *RectCollider) Collisions(oth Polygon) bool {
	return PolygonsIntersect(r.polygon(), oth)
}

// PolygonsIntersect retourne vrai si les deux polygones spécifiés se chevauchent.
// Sinon, retourne faux.
func PolygonsIntersect(a, b Polygon) bool {
	for _, poly := range [2]Polygon{a, b} {
		for i := 0; i < len(poly.points); i++ {
			j := (i + 1) % len(poly.points)

			p1 := poly.points[i]
			p2 := poly.points[j]

			normal := &Point{X: p2.Y - p1.Y, Y: p1.X - p2.X}

			minA, maxA := float32(math.MaxFloat32), float32(-math.MaxFloat32)
			for _, point := range a.points {
				projected := normal.X*point.X + normal.Y*point.Y
				if projected < minA {
					minA = projected
				}
				if projected > maxA {
					maxA = projected
				}
			}

			minB, maxB := float32(math.MaxFloat32), float32(-math.MaxFloat32)
			for _, point := range b.points {
				projected := normal.X*point.X + normal.Y*point.Y
				if projected < minB {
					minB = projected
				}
				if projected > maxB {
					maxB = projected
				}
			}

			if maxA < minB || maxB < minA {
				return false
			}
		}
	}
	return true
}
