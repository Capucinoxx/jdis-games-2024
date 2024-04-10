package manager

import (
	"log"
	"sync"
	"time"

	"github.com/capucinoxx/forlorn/model"
	"github.com/capucinoxx/forlorn/utils"
)

// tickrate est le nombre de mises à jour du jeu par seconde.
const tickrate = 3

// GameManager est responsable de la gestion de l'état du jeu, des joueurs et
// de la boucle de jeu. Il est également responsable de la gestion des messages
// entrants des joueurs et de la synchronisation des mises à jour du jeu.
type GameManager struct {
	tickStart time.Time
	am        *AuthManager
	nm        *NetworkManager
	state     *model.GameState
	mu        sync.Mutex
}

// NewGameManager crée un nouveau gestionnaire de jeu avec le serveur de jeu et
// le gestionnaire de réseau spécifiés.
func NewGameManager(am *AuthManager, nm *NetworkManager) *GameManager {
	return &GameManager{
		state: model.NewGameState(),
		am:    am,
		nm:    nm,
	}
}

// RegisterPlayer ajoute un joueur à l'état du jeu et à la liste des clients.
func (gm *GameManager) RegisterPlayer(conn model.Connection) {
	// players := gm.state.Players()

	spawn := []float32{0.0, 0.0} // TODO: Generate spawn position
	player := model.NewPlayer(0, spawn[0], spawn[1], conn)

	gm.nm.Register(player)
	gm.state.AddPlayer(player)
	if gm.state.InProgess() {
		// TODO: send game start to player
	}
}

// UnregisterPlayer supprime un joueur de l'état du jeu et de la liste des clients.
func (gm *GameManager) UnregisterPlayer(conn model.Connection) {
	players := gm.state.Players()

	for _, p := range players {
		if p.Client.Connection == conn {
			gm.state.RemovePlayer(p)
			break
		}
	}
}

// Init initialise le gestionnaire de jeu. Il démarre le serveur de jeu et
// attend les connexions des joueurs.
func (gm *GameManager) Init() error {
	return gm.nm.Start()
}

// Start démarre le gestionnaire de jeu. Il initialise le serveur de jeu et
// démarre la boucle de jeu.
func (gm *GameManager) Start() {
	gm.state.Start()
	go gm.gameLoop()
}

// State retourne l'état actuel du jeu.
func (gm *GameManager) State() (model.Map, int) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	state := gm.state
	return *state.Map, 0
}

// process traite les messages entrants des joueurs. Il met à jour l'état du jeu
// en fonction des messages entrants. Il met également à jour les informations
// des joueurs en fonction des messages entrants. La méthode prend en charge
// l'authentification des joueurs et la mise à jour des informations des joueurs.
func (gm *GameManager) process(p *model.Player, players []*model.Player, timestep float32) {
	for len(p.Client.In) != 0 {
		message := <-p.Client.In
		utils.Log("game", "process", "Player %d received message %v", p.ID, message)

		switch msgType := message.MessageType; msgType {
		case model.Spawn:
			// Lorsqu'un joueur se connecte, il doit envoyer un message de type Spawn
			// pour s'authentifier. Si le jeton d'authentification est valide, le joueur
			// est autorisé à rejoindre la partie. Sinon, le joueur est déconnecté.
			tkn := message.Body.(string)
			if !gm.am.Authenticate(tkn) {
				gm.nm.ForceDisconnect(p)
				continue
			}
			log.Printf("Player %d spawned", p.ID)
		case model.Position:
			// Lorsqu'un joueur envoie un message de type Position, cela signifie
			// qu'il a bougé ou tourné. Le message contient les nouvelles coordonnées
			// du joueur. Ces coordonnées sont utilisées pour mettre à jour la position
			// du joueur.
			p.Controls = message.Body.(model.Controls)
			p.Update(players, gm.state, timestep)

		}
	}
}

// gameLoop est la boucle principale du jeu. Il gère les mises à jour du jeu,
// les entrées des joueurs et les sorties des joueurs. La boucle de jeu appelle
// cette méthode dans une goroutine séparée. Cela permet à la boucle de jeu de
// continuer à s'exécuter même si le jeu est occupé. La boucle de jeu est
// responsable de la synchronisation des mises à jour du jeu.
func (gm *GameManager) gameLoop() {
	interval := time.Duration((int(1000 / tickrate))) * time.Millisecond
	timestep := float32(interval/time.Millisecond) / 1000.0

	ticker := time.NewTicker(interval)
	gm.nm.BroadcastGameStart(gm.state)

	for range ticker.C {
		gm.tickStart = time.Now()
		players := gm.state.Players()

		for _, p := range players {
			gm.process(p, players, timestep)
			// handle respawn
		}

		gm.nm.BroadcastGameState(gm.state)
	}
	ticker.Stop()

	gm.nm.BroadcastGameEnd()
	time.Sleep(10 * time.Second)
	gm.Start()
}
