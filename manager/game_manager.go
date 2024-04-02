package manager

import (
	"time"

	"github.com/capucinoxx/forlorn/model"
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

// HeartBeat envoie un message de battement de coeur au serveur de jeu toutes les 5 secondes.
// Cela permet au serveur de jeu de savoir que le gestionnaire de jeu est toujours actif.
func (gm *GameManager) HeartBeat() {
	time.Sleep(time.Second)

	ticker := time.NewTicker(time.Second * 5)
	// gm.server.Init(gm.nm.Address())
	for range ticker.C {
		// gm.server.SendHeartbeat(gm.state)
	}
}

// Start démarre le gestionnaire de jeu. Il initialise le serveur de jeu et
// démarre la boucle de jeu.
func (gm *GameManager) Start() {
	gm.state.Start()
	go gm.gameLoop()
}

// process traite les messages entrants des joueurs. Il met à jour l'état du jeu
// en fonction des messages entrants. Il met également à jour les informations
// des joueurs en fonction des messages entrants. La méthode prend en charge
// l'authentification des joueurs et la mise à jour des informations des joueurs.
func (gm *GameManager) process(p *model.Player, players []*model.Player, timestep float32) {
	for len(p.Client.In) != 0 {
		message := <-p.Client.In
		switch msgType := message.MessageType; msgType {
		case 0:
			// token := message.Body.(string)
			// TODO: authenticate player
		case 1:
			// TODO: update player information
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
	gm.nm.BroadcastGameState()

	for range ticker.C {
		gm.tickStart = time.Now()
		players := gm.state.Players()

		for _, p := range players {
			gm.process(p, players, timestep)
			// handle respawn
		}

		gm.nm.BroadcastGameState()

	}
	ticker.Stop()

	// gm.server.EndGame(gm.state)
	time.Sleep(10 * time.Second)
	gm.Start()
}
