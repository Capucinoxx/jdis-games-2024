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
	vertices []*Point
}

// String retourne une représentation en chaîne de caractères du polygone.
func (p Polygon) String() string {
	str := "["
	for _, point := range p.vertices {
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

// polygon retourne le polygone représenté par le RectCollider.
func (r *RectCollider) polygon() Polygon {
	return Polygon{
		vertices: []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d},
	}
}

// Collisions retourne vrai si le RectCollider entre en collision avec
// le polygone spécifié. Sinon retourne faux.
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

// PolygonsIntersect retourne vrai si les deux polygones spécifiés se chevauchent.
// Sinon, retourne faux.
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
