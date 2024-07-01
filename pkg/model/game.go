package model

import (
	"sync"
	"time"

	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/config"
)

// GameState représente l'état actuel du jeu.
type GameState struct {
	startTime    time.Time
	inProgress   bool
	lastPlayerID int64

	playerCount  int
	players      map[string]*Player
	coins       []*Scorer

  Map          Map
  spawns      []*Point
  spawnIndex  int
	mu           *sync.RWMutex
}

// NewGameState crée un nouvel état de jeu.
func NewGameState(m Map) *GameState {
  return &GameState{
    spawns:       []*Point{},
    spawnIndex:   0,
    inProgress:   false,
    lastPlayerID: 0,
    coins:        []*Scorer{},
    playerCount:  0,
    players:      make(map[string]*Player),
    Map:          m,
    mu:           &sync.RWMutex{},
  }
}

func (gs *GameState) GetSpawnPoint() *Point {
  
  spawn := gs.spawns[gs.spawnIndex]
  gs.spawnIndex = (gs.spawnIndex + 1) % len(gs.spawns)
  return spawn
}

func (gs *GameState) SetSpawns(spawns []*Point) {
  gs.spawnIndex = 0
  gs.spawns = spawns
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
    if p.Client.Connection.Identifier() != "" {
		  players = append(players, p)
    }
	}
	return players
}

func (gs *GameState) Coins() []*Scorer {
  gs.mu.RLock()
  defer gs.mu.RUnlock()
  return gs.coins
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

	gs.players[p.Nickname] = p
	gs.playerCount++
	return gs.playerCount
}

// RemovePlayer supprime un joueur de l'état du jeu.
func (gs *GameState) RemovePlayer(p *Player) int {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if _, ok := gs.players[p.Nickname]; !ok {
		return gs.playerCount
	}

	delete(gs.players, p.Nickname)
	gs.playerCount--
	return gs.playerCount
}

// Start démarre le jeu. Inialise les joueurs,
// démarre le chronomètre et met un drapeau pour indiquer
// que le jeu est en cours.
func (gs *GameState) Start() {
	gs.Map.Setup()
  gs.SetSpawns(gs.Map.Spawns(0))

  players := gs.Players()
  for _, p := range players {
    p.Respawn(gs)
  }

  gs.coins = make([]*Scorer, 0, config.NumCoins)
  for i := 0; i < config.NumCoins; i++ {
    gs.coins = append(gs.coins, NewCoin())
  }

	gs.startTime = time.Now()

	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.inProgress = true
}

type GameMessage struct {
  Players []*Player
  Coins []*Scorer
}

type GameInfo struct {
  Players []PlayerInfo
  Coins []struct{
    Uuid [16]byte
    Value int32
    Pos Point
  }
}


func (gs *GameMessage) Encode(w codec.Writer) (err error) {
  if err = w.WriteInt32(int32(len(gs.Players))); err != nil {
    return
  }

  for _, p := range gs.Players {
    if err = p.Encode(w); err != nil {
      return
    }
  }

  if err = w.WriteInt32(int32(len(gs.Coins))); err != nil {
    return
  }

  for _, c := range gs.Coins {
    if err = c.Encode(w); err != nil {
      return
    }
  }

  return
}

func (g *GameInfo) Decode(r codec.Reader) (err error) {
  var size int32
  if size, err = r.ReadInt32(); err != nil {
    return
  }

  g.Players = make([]PlayerInfo, 0, size)
  for i := int32(0); i < size; i++ {
    p := PlayerInfo{}
    if err = p.Decode(r); err != nil {
      return
    }
    g.Players = append(g.Players, p)
  }

  if size, err = r.ReadInt32(); err != nil {
    return
  }
  
  g.Coins = make([]struct{
    Uuid [16]byte
    Value int32
    Pos Point
  }, 0, size)
  for i := int32(0); i < size; i++ {
    c := struct{
      Uuid [16]byte
      Value int32
      Pos Point
    }{}

    var id []byte
    if id, err = r.ReadBytes(16); err != nil {
      return
    }
    copy(c.Uuid[:], id)
    if err = c.Pos.Decode(r); err != nil {
      return
    }
    if c.Value, err = r.ReadInt32(); err != nil {
      return
    }

    g.Coins = append(g.Coins, c)
  }

  return
}
