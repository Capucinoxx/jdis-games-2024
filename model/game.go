package model

import (
	"sync"
	"time"
)

const gameLength = 10

// GameState représente l'état actuel du jeu.
type GameState struct {
	startTime    time.Time
	length       time.Duration
	inProgress   bool
	lastPlayerID int64
	playerCount  int
	players      map[uint8]*Player
	Map          *Map
	mu           *sync.RWMutex
}

// NewGameState crée un nouvel état de jeu.
func NewGameState() *GameState {
	return &GameState{
		length:       gameLength * time.Minute,
		inProgress:   false,
		lastPlayerID: 0,
		playerCount:  0,
		players:      make(map[uint8]*Player),
		Map:          &Map{},
		mu:           &sync.RWMutex{},
	}
}

// InProgess retourne vrai si le jeu est en cours.
// Sinon, il retourne faux.
func (gs *GameState) InProgess() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.inProgress
}

// Players retourne une liste de tous les joueurs.
func (gd *GameState) Players() []*Player {
	gd.mu.RLock()
	defer gd.mu.RUnlock()

	players := make([]*Player, 0, len(gd.players))
	for _, p := range gd.players {
		players = append(players, p)
	}
	return players
}

// PlayerCount retourne le nombre de joueurs.
func (gs *GameState) PlayerCount() int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.playerCount
}

// AddPlayer ajoute un joueur à l'état du jeu.
func (gs *GameState) AddPlayer(p *Player) int {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.players[p.ID] = p
	gs.playerCount++
	return gs.playerCount
}

// RemovePlayer supprime un joueur de l'état du jeu.
func (gs *GameState) RemovePlayer(p *Player) int {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if _, ok := gs.players[p.ID]; !ok {
		return gs.playerCount
	}

	delete(gs.players, p.ID)
	gs.playerCount--
	return gs.playerCount
}

// Start démarre le jeu. Inialise les joueurs,
// démarre le chronomètre et met un drapeau pour indiquer
// que le jeu est en cours.
func (gs *GameState) Start() {
	// players := gs.Players()

	// TODO: initialize player informations (health, score, position)

	gs.startTime = time.Now()

	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.inProgress = true
}
