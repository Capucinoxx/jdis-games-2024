package model

import (
	"sync"
	"time"
)

type GameState struct {
	startTime  time.Time
	inProgress bool
	freeze     bool

	players     map[string]*Player
	cachedScore map[string]int
	coins       *Scorers

	Map        Map
	spawns     []*Point
	spawnIndex int
	mu         *sync.RWMutex
}

func NewGameState(m Map) *GameState {
	return &GameState{
		spawns:      []*Point{},
		spawnIndex:  0,
		inProgress:  false,
		freeze:      false,
		coins:       NewScorers(),
		players:     make(map[string]*Player),
		cachedScore: make(map[string]int),
		Map:         m,
		mu:          &sync.RWMutex{},
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

func (gs *GameState) SetCoins(coins []*Scorer) {
	gs.coins.Add(coins...)
}

func (gs *GameState) SetFreeze(b bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	gs.freeze = b
}

func (gs *GameState) IsFreeze() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.freeze
}

func (gs *GameState) InProgess() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.inProgress
}

func (gd *GameState) Players() []*Player {
	gd.mu.RLock()
	defer gd.mu.RUnlock()

	players := make([]*Player, 0, len(gd.players))
	for _, p := range gd.players {
		if p.Client.GetConnection().Identifier() != "" {
			players = append(players, p)
		}
	}
	return players
}

func (gs *GameState) PlayersScore() []PlayerScore {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	players := make([]PlayerScore, 0, len(gs.players))
	for _, p := range gs.players {
		identifier := p.Client.GetConnection().Identifier()
		if identifier != "" {
			if v, ok := gs.cachedScore[identifier]; ok {
				players = append(players, PlayerScore{Name: p.Nickname, Score: p.Score() - v})
			} else {
				players = append(players, PlayerScore{Name: p.Nickname, Score: p.Score()})
			}
			gs.cachedScore[identifier] = p.score
		}
	}
	return players
}

func (gs *GameState) Coins() *Scorers {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.coins
}

func (gs *GameState) AddPlayer(username string, color int, conn Connection) *Player {
	var player *Player
	var ok bool

	gs.mu.Lock()
	player, ok = gs.players[username]
	gs.mu.Unlock()

	if ok {
		player.Client.SetConnection(conn)
		player.Client.In = make(chan ClientMessage, 10)
		player.Client.Out = make(chan []byte, 10)
		return player
	}

	spawn := &Point{0, 0}
	if gs.InProgess() {
		spawn = gs.GetSpawnPoint()
	}
	player = NewPlayer(username, color, spawn, conn)
	gs.mu.Lock()
	gs.players[username] = player
	gs.mu.Unlock()
	return player
}

func (gs *GameState) RemovePlayer(p *Player) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	delete(gs.players, p.Nickname)
}

func (gs *GameState) Start() {
	gs.mu.Lock()
	if gs.inProgress {
		gs.mu.Unlock()
		return
	}
	gs.mu.Unlock()

	gs.Map.Setup()
	gs.SetSpawns(gs.Map.Spawns(0))

	gs.startTime = time.Now()

	gs.mu.Lock()
	defer gs.mu.Unlock()
	for _, p := range gs.players {
		p.score = 0
	}
	clear(gs.cachedScore)
	gs.inProgress = true
}

func (gs *GameState) Reset(scorers []*Scorer) {
	players := gs.Players()
	for _, p := range players {
		p.Respawn(gs)
	}

	gs.coins.Set(scorers)
}

func (gs *GameState) Stop() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.inProgress = false
}
