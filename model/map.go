package model

import "math"

// Point représente un point "continu" dans un espace 2D.
type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// DirectionTo retourne un vecteur normalisé pointant vers la destination à
// partir du point actuel.
func (p *Point) DirectionTo(dest *Point) *Point {
	dir := &Point{X: 0, Y: 0}

	dir.X = p.X - dest.X
	dir.Y = p.Y - dest.Y
	dir.normalize()

	return dir
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

// normalize normalise le vecteur. Cela signifie que la longueur du vecteur est
// égale à 1.
func (p *Point) normalize() {
	length := math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
	p.X /= float32(length)
	p.Y /= float32(length)
}

type ColliderType uint8

// Collider représente un polygone qui représentant une collision dans le jeu.
type Collider struct {
	Points []*Point     `json:"points"`
	Type   ColliderType `json:"type"`
}

// Map représente une carte de jeu. Elle contient des informations sur les
// collisions et les points de spawn.
type Map struct {
	Collider []*Collider
	Spawns   []*Point
}
