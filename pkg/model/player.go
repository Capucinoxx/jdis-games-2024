package model

import (
	"fmt"
	"time"
)

// Connection représente une connexion réseau. Elle peut être utilisée pour lire
// et écrire des données sur le réseau. Elle peut également être utilisée pour
// identifier la connexion. La connexion doit être fermée après utilisation.
type Connection interface {
	Identifier() string
	Close(time.Duration)

	PrepareRead(int64, time.Duration)
	Read() ([]byte, error)

	PrepareWrite(time.Duration)
	Write([]byte) error

	Ping(time.Duration)
}

const (
	// playerSize est la taille du joueur.
	playerSize = 0.1

	// defaultHealth est la vie par défaut du joueur.
	defaultHealth = 100
)

// Controls représente les contrôles du joueur.
// orsqu'un contrôle est activé, le joueur effectue
// cette action.
type Controls struct {
	Rotation uint32
}

type Player struct {
	ID               uint8
	Token            string
	Nickname         string
	Health           int
	Score            float64
	respawnCountdown float32
	Client           *Client
	Controls         Controls
	Collider         *RectCollider
}

func NewPlayer(id uint8, x float32, y float32, conn Connection) *Player {
	return &Player{
		ID:       id,
		Collider: NewRectCollider(x, y, playerSize),
		Health:   defaultHealth,
		Client: &Client{
			Out:        make(chan []byte, 10),
			In:         make(chan ClientMessage, 10),
			Connection: conn,
		},
	}
}

// IsAlive retourne vrai si le joueur est en vie. Sinon, il retourne faux.
// Un joueur est considéré comme mort si sa vie est inférieure ou égale à 0.
func (p *Player) IsAlive() bool {
	return p.Health > 0
}

// Update met à jour l'état du joueur en fonction de l'état actuel du jeu.
func (p *Player) Update(players []*Player, game *GameState, dt float32) {
	m := game.Map
	if !p.IsAlive() {
		p.respawnCountdown += dt
		return
	}

	p.HandleMovement(players, m, dt)
}

// HandleMovement gère le mouvement du joueur en fonction de ses contrôles.
func (p *Player) HandleMovement(players []*Player, m *Map, dt float32) {
	r := p.Collider

	hasCollision := p.checkCollisionWithPlayers(players) || p.checkCollisionWithMap(m)

	p.updateVelocity(dt, hasCollision)
	p.updateRotation()
	if !hasCollision {
		p.applyMovement()
	}

	r.Rotation = (r.Rotation + p.Controls.Rotation) % 360
}

// String retourne une représentation en chaîne de caractères du joueur.
func (p *Player) String() string {
	return fmt.Sprintf("[%d: { pos: (%f, %f), v: %f, rot: %d, health: %d }]", p.ID, p.Collider.Pivot.X, p.Collider.Pivot.Y, p.Collider.velocity, p.Collider.Rotation, p.Health)
}

// checkCollisionWithPlayers retourne vrai si le joueur entre en collision avec
// un autre joueur. Sinon, retourne faux.
func (p *Player) checkCollisionWithPlayers(players []*Player) bool {
	for _, ennemy := range players {
		if ennemy.ID == p.ID || !ennemy.IsAlive() {
			continue
		}

		if p.Collider.Collisions(ennemy.Collider.polygon()) {
			return true
		}
	}

	return false
}

// checkCollisionWithMap retourne vrai si le joueur entre en collision avec la carte.
// Sinon, retourne faux.
func (p *Player) checkCollisionWithMap(m *Map) bool {
	for _, collider := range m.Colliders {
		if p.Collider.Collisions(collider.polygon()) {
			fmt.Println(p.Collider.polygon().String())
			fmt.Println(collider.polygon().String())
			return true
		}
	}

	return false
}

// updateVelocity met à jour la vitesse du joueur en fonction de la collision.
// Si le joueur entre en collision, sa vitesse est réduite à 0.
func (p *Player) updateVelocity(dt float32, hasCollision bool) {
	r := p.Collider
	r.velocity = defaultForwardSpeed
}

// updateRotation met à jour la rotation du joueur en fonction de ses contrôles.
func (p *Player) updateRotation() {
	p.applyRotation(p.Controls.Rotation)
}

// applyRotation applique la rotation spécifiée au joueur. La rotation est appliquée
// à tous les points du joueur.
func (p *Player) applyRotation(rd uint32) {
	r := p.Collider
	points := []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d, r.look}
	for _, point := range points {
		r.rotate(rd, point)
	}
	r.CalculDirection()
}

// applyMovement applique le mouvement du joueur. Le mouvement est appliqué à tous
// les points du joueur.
func (p *Player) applyMovement() {
	r := p.Collider
	points := []*Point{r.rect.a, r.rect.b, r.rect.c, r.rect.d, r.look, r.Pivot}
	for _, point := range points {
		point.X += r.dir.X * r.velocity
		point.Y += r.dir.Y * r.velocity
	}
}