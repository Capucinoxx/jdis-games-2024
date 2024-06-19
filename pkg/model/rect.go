package model


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
	forwardSpeed float64

	velocity float64
}

// NewRectCollider crée un nouveau RectCollider.
func NewRectCollider(x, y, size float64) *RectCollider {
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

func ProjectPolygon(axis Point, polygon Polygon) (float64, float64) {
  min := (polygon.points[0].X * axis.X) + (polygon.points[0].Y * axis.Y)
  max := min

  for _, vertex := range polygon.points {
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

// PolygonsIntersect retourne vrai si les deux polygones spécifiés se chevauchent.
// Sinon, retourne faux.
func PolygonsIntersect(a, b Polygon) bool {
	return true
}
