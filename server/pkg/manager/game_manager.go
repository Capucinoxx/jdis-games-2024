package manager

// GameManager orchestrates the game logic and state management for a multiplayer game.This package provides
// the necessary tools to manage player connections, game state updates, and game loop execution. It integrates
// tightly with the AuthManager, NetworkManager, and RoundManager to provide a seamless gaming experience.
//
// The GameManager is responsible for:
// - Registering and unregistering player connections, including handling spectators and authenticated players.
// - Initializing and starting the game server.
// - Managing the game loop, which includes processing player actions and updating the game state.
// - Broadcasting game state updates, game start, and game end messages to all connected clients.
// - Handling player-specific actions such as respawn and damage.
//
// The GameManager utilizes a mutex to ensure thread safety when accessing and modifying the game state.
// It also abstracts the details of authentication, network communication, and round management via interfaces.
//
// Usage of this package involves creating an instance of GameManager with specific instances of AuthManager,
// NetworkManager, and RoundManager, along with the initial game map. The GameManager then handles the game
// lifecycle, including player management and game state updates.

import (
	"fmt"
	"time"

	"github.com/capucinoxx/jdis-games-2024/consts"
	"github.com/capucinoxx/jdis-games-2024/pkg/model"
	"github.com/capucinoxx/jdis-games-2024/pkg/utils"
)

// RoundManager is an interface for managing game rounds and tick.
type RoundManager interface {
	Restart()
	Tick()
	CurrentTick() int
	CurrentRound() int8
	SetState(*model.GameState)
	HasEnded() bool
}

// GameManager maintains the game state and manages the game loop.
type GameManager struct {
	tickStart time.Time
	am        *AuthManager
	nm        *NetworkManager
	rm        RoundManager
	sm        *ScoreManager
	state     *model.GameState
}

// NewGameManager creates a new GameManager with the specified authentication, network, and round managers, and initial
// game map
func NewGameManager(am *AuthManager, nm *NetworkManager, rm RoundManager, sm *ScoreManager, m model.Map) *GameManager {
	state := model.NewGameState(m)
	rm.SetState(state)

	return &GameManager{
		state: state,
		am:    am,
		nm:    nm,
		sm:    sm,
		rm:    rm,
	}
}

// RegisterConnection registers a new connection, either as a player or a spectator.
func (gm *GameManager) RegisterConnection(conn model.Connection, adminToken string) error {
	if conn.Identifier() == "" {
		gm.addSpectator(conn, adminToken)
		return nil
	} else {
		return gm.addPlayer(conn)
	}
}

// addSpectator adds a new spectator to the game. A spectator is a client that is not authenticated as a player.
// Spectators receive game state updates but cannot interact with the game.
func (gm *GameManager) addSpectator(conn model.Connection, token string) {
	client := &model.Client{
		Out: make(chan []byte, 10),
	}
	client.SetConnection(conn)

	isAdmin := false
	if token != "" {
		_, _, isAdmin, _ = gm.am.Authenticate(token)
		conn.SetAdmin(isAdmin)
	}

	gm.nm.Register(client)
	if gm.state.InProgess() {
		gm.nm.Send(client, gm.nm.protocol.Encode(&model.ClientMessage{
			MessageType: model.MessageMapState,
			Body: model.MessageMapStateToEncode{
				Map:     gm.state.Map,
				IsAdmin: isAdmin,
			},
		}))
	}
}

// addPlayer adds a new player to the game. A player is a client that is authenticated and can interact with the game.
func (gm *GameManager) addPlayer(conn model.Connection) error {
	username, color, isAdmin, ok := gm.am.Authenticate(conn.Identifier())
	if !ok {
		return fmt.Errorf("unknown token")
	}
	conn.SetAdmin(isAdmin)

	player := gm.state.AddPlayer(username, color, conn)
	gm.nm.Register(player.Client)

	if gm.state.InProgess() {
		gm.nm.Send(player.Client, gm.nm.protocol.Encode(&model.ClientMessage{
			MessageType: model.MessageMapState,
			Body: model.MessageMapStateToEncode{
				Map:     gm.state.Map,
				IsAdmin: isAdmin,
				Storage: player.Storage(),
			},
		}))
	}

	return nil
}

func (gm *GameManager) RemoveConnection(conn model.Connection) {}

// Initialize starts the network manager and prepares the game for execution.
func (gm *GameManager) Initialize() error {
	return gm.nm.Start()
}

func (gm *GameManager) Freeze(b bool) {
	gm.state.SetFreeze(b)
}

// Start starts the game, initializing the game state and starting the game loop.
func (gm *GameManager) Start() {
	if !gm.state.IsFreeze() && !gm.state.InProgess() {
		gm.state.Start()

		gm.rm.Restart()
		go gm.gameLoop()
	}
}

// Kill foribly removes a player from the game by setting their health to 0.
// This is used for debugging purposes.
func (gm *GameManager) Kill(name string) {
	for _, player := range gm.state.Players() {
		if player.Nickname == name {
			player.TakeDmg(1_000_000)
			return
		}
	}
}

// process processes player actions and updates the game state.
func (gm *GameManager) process(p *model.Player, players []*model.Player, timestep float64, handleAction bool) {
	for len(p.Client.In) != 0 {
		message := <-p.Client.In

		switch msgType := message.MessageType; msgType {
		case model.MessagePlayerAction:
			if handleAction {
				p.Controls = message.Body.(model.Controls)
			}
		}
	}

	p.Update(players, gm.state, timestep)
}

// gameLoop is the main game loop that handles game state updates and broadcasting game state to clients.
func (gm *GameManager) gameLoop() {
	interval := time.Duration((int(1000 / consts.Tickrate))) * time.Millisecond
	timestep := float64(interval/time.Millisecond) / 1000.0

	for _, p := range gm.state.Players() {
		p.ClearStorage()
	}

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

		ok := gm.state.Coins().Update()
		if ok {
			gm.state.Stop()
			break
		}

		if count == 10 {
			gm.nm.BroadcastGameState(gm.state, int32(gm.rm.CurrentTick()), gm.rm.CurrentRound())
			scores := gm.state.PlayersScore()
			count = 0
			gm.sm.Adds(scores)
		}

		count++
		if gm.rm.HasEnded() {
			gm.state.Stop()
			break
		}
	}
	ticker.Stop()

	gm.nm.BroadcastGameEnd()
	go func() {
		if err := gm.sm.Persist(); err != nil {
			utils.Log("error", "persist", "mongo persistance error %s", err)
		}
	}()
	gm.Start()
}
