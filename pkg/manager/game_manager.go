package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/capucinoxx/forlorn/pkg/model"
)

// tickrate est le nombre de mises à jour du jeu par seconde.
const tickrate = 30

type RoundManager interface {
	Restart()
	Tick()
  CurrentTick() int
  CurrentRound() int8
  SetState(*model.GameState)
	HasEnded() bool
}

// GameManager est responsable de la gestion de l'état du jeu, des joueurs et
// de la boucle de jeu. Il est également responsable de la gestion des messages
// entrants des joueurs et de la synchronisation des mises à jour du jeu.
type GameManager struct {
	tickStart time.Time
	am        *AuthManager
	nm        *NetworkManager
	rm        RoundManager
	state     *model.GameState
	mu        sync.Mutex
}

// NewGameManager crée un nouveau gestionnaire de jeu avec le serveur de jeu et
// le gestionnaire de réseau spécifiés.
func NewGameManager(am *AuthManager, nm *NetworkManager, rm RoundManager, m model.Map) *GameManager {
  state := model.NewGameState(m)
  rm.SetState(state)

	return &GameManager{
		state: state,
		am:    am,
		nm:    nm,
		rm:    rm,
	}
}


func (gm *GameManager) Register(conn model.Connection) error {
  if (conn.Identifier() == "") {
    gm.RegisterSpectator(conn)
    return nil
  } else {
     return gm.RegisterPlayer(conn)
  }
}

func (gm *GameManager) RegisterSpectator(conn model.Connection) {
  client := &model.Client{
    Out: make(chan []byte, 10),
    Connection: conn,
  }
  gm.nm.Register(client)
  if gm.state.InProgess() {
    gm.nm.Send(client, gm.nm.protocol.Encode(0, 0, &model.ClientMessage{
      MessageType: model.GameStart,
      Body: gm.state.Map,
    }))
  }
}

// RegisterPlayer ajoute un joueur à l'état du jeu et à la liste des clients.
func (gm *GameManager) RegisterPlayer(conn model.Connection) error {
  username, color, ok := gm.am.Authenticate(conn.Identifier())
  if !ok {
    return fmt.Errorf("Unknown token")
  }

  spawn := &model.Point{X: 0, Y: 0}

  if gm.state.InProgess() {
    spawn = gm.state.GetSpawnPoint()
  }
	player := model.NewPlayer(username, color, spawn, conn)
  
	gm.nm.Register(player.Client)
	gm.state.AddPlayer(player)


	if gm.state.InProgess() {
    gm.nm.Send(player.Client, gm.nm.protocol.Encode(0, 0, &model.ClientMessage{
      MessageType: model.GameStart,
      Body: gm.state.Map,
    }))
	}

  return nil
}

// UnregisterPlayer supprime un joueur de l'état du jeu et de la liste des clients.
func (gm *GameManager) Unregister(conn model.Connection) {}

// Init initialise le gestionnaire de jeu. Il démarre le serveur de jeu et
// attend les connexions des joueurs.
func (gm *GameManager) Init() error {
	return gm.nm.Start()
}

// Start démarre le gestionnaire de jeu. Il initialise le serveur de jeu et
// démarre la boucle de jeu.
func (gm *GameManager) Start() {
	gm.state.Start()

  gm.rm.Restart()
	go gm.gameLoop()
}

// State retourne l'état actuel du jeu.
func (gm *GameManager) State() (model.Map, int) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	state := gm.state
	return state.Map, 0
}

func (gm *GameManager) Kill(name string) {
  for _, player := range gm.state.Players() {
    if player.Nickname == name {
      player.TakeDmg(1_000_000)
      return
    }
  }
}


// process traite les messages entrants des joueurs. Il met à jour l'état du jeu
// en fonction des messages entrants. Il met également à jour les informations
// des joueurs en fonction des messages entrants. La méthode prend en charge
// l'authentification des joueurs et la mise à jour des informations des joueurs.
func (gm *GameManager) process(p *model.Player, players []*model.Player, timestep float64, handleAction bool) {
	for len(p.Client.In) != 0 {
		message := <-p.Client.In

		switch msgType := message.MessageType; msgType {
		case model.Spawn:
			// Lorsqu'un joueur se connecte, il doit envoyer un message de type Spawn
			// pour s'authentifier. Si le jeton d'authentification est valide, le joueur
			// est autorisé à rejoindre la partie. Sinon, le joueur est déconnecté.
			tkn := message.Body.(string)
      if user, color, ok := gm.am.Authenticate(tkn); !ok {
				gm.nm.ForceDisconnect(p.Client.Connection)
				continue
			} else {
        p.Nickname = user
        p.Color = color
      }

		case model.Action:
      if handleAction {
			  p.Controls = message.Body.(model.Controls)
      }
      break;
		}
	}

  p.Update(players, gm.state, timestep)
}

// gameLoop est la boucle principale du jeu. Il gère les mises à jour du jeu,
// les entrées des joueurs et les sorties des joueurs. La boucle de jeu appelle
// cette méthode dans une goroutine séparée. Cela permet à la boucle de jeu de
// continuer à s'exécuter même si le jeu est occupé. La boucle de jeu est
// responsable de la synchronisation des mises à jour du jeu.
func (gm *GameManager) gameLoop() {
	interval := time.Duration((int(1000 / tickrate))) * time.Millisecond
	timestep := float64(interval/time.Millisecond) / 1000.0

	ticker := time.NewTicker(interval)
	gm.nm.BroadcastGameStart(gm.state)

  count := 0 
	for range ticker.C {
		gm.tickStart = time.Now()
		players := gm.state.Players()

		gm.rm.Tick()

		for _, p := range players {
			gm.process(p, players, timestep, true)
			p.HandleRespawn(gm.state)
		}
    
    if count == 10 {
		  gm.nm.BroadcastGameState(gm.state, int32(gm.rm.CurrentTick()), gm.rm.CurrentRound())
      count = 0
    }
      
    count++
		if gm.rm.HasEnded() {
      gm.state.Stop()
			break
		}
	}
	ticker.Stop()

	gm.nm.BroadcastGameEnd()
  
  
	time.Sleep(10 * time.Second)
	gm.Start()
}
