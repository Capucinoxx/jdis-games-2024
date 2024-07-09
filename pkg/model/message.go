package model

import (
	"github.com/capucinoxx/forlorn/pkg/codec"
)

type MessageType uint8

const (
	// MessageGameState is used when the server wants to send the dynamic state of the game
	// (state of the players and coins
	// Encode: MessageGameStateToEncode.Encode()
	// Decode: MessageGameStateToDecode.Decode()
	//
	//
	// +-------------------+------------------------------------------+
	// |          Binary Representation                               |
	// +-------------------+------------------------------------------+
	// | Field             | Description                              |
	// +-------------------+------------------------------------------+
	// | 4 bytes (int32)   | current tick                             |
	// | 1 byte  (int8)    | current round (0/1)                      |
	// | 4 bytes (int32)   | number of players                        |
	// +-------------------+------------------------------------------+
	// | For each player (0 .. number of players) do                  |
	// +-------------------+------------------------------------------+
	// | n bytes (string)  | player name (read until \0)              |
	// | 4 bytes (int32)   | player color                             |
	// | 4 bytes (int32)   | player health                            |
	// | 8 bytes (float64) | player x axis position                   |
	// | 8 bytes (float64) | player y axis position                   |
	// | 1 byte  (bool)    | if player has destination (0/1)          |
	// +-------------------+------------------------------------------+
	// | If player has destination:                                   |
	// +-------------------+------------------------------------------+
	// | 8 bytes (float64) | player x axis destination                |
	// | 8 bytes (float64) | player y axis destination                |
	// +-------------------+------------------------------------------+
	// | End if player has destination                                |
	// +-------------------+------------------------------------------+
	// | 4 bytes (int32)   | number of player projectiles             |
	// +-------------------+------------------------------------------+
	// | For each player projectile (0 .. number of projectiles) do   |
	// +-------------------+------------------------------------------+
	// | 16 bytes (string) | projectile unique id                     |
	// | 8 bytes (float64) | projectile x axis position               |
	// | 8 bytes (float64) | projectile y axis position               |
	// | 8 bytes (float64) | projectile x axis destination            |
	// | 8 bytes (float64) | projectile y axis destination            |
	// +-------------------+------------------------------------------+
	// | End for each player projectile                               |
	// +-------------------+------------------------------------------+
	// | 8 bytes (float64) | blade x axis start position              |
	// | 8 bytes (float64) | blade y axis start position              |
	// | 8 bytes (float64) | blade x axis end position                |
	// | 8 bytes (float64) | blade y axis end position                |
	// +-------------------+------------------------------------------+
	// | End for each player                                          |
	// +--------------------------------------------------------------+
	MessageGameState = 1

	MessagePlayerAction = 3

	// +-------------------+------------------------------------------+
	// |          Binary Representation                               |
	// +-------------------+------------------------------------------+
	// | Field             | Description                              |
	// +-------------------+------------------------------------------+
	// | 1 byte  (int8)    | size of discrete map (size x size)       |
	// +-------------------+------------------------------------------+
	// | For each cell in discrete map (size * size times)            |
	// +-------------------+------------------------------------------+
	// | 1 byte  (uint8)   | number of walls in cell                  |
	// +-------------------+------------------------------------------+
	// | End for each cell in discrete map                            |
	// +-------------------+------------------------------------------+
	// | 4 bytes (int32)   | number of walls                          |
	// +-------------------+------------------------------------------+
	// | For each wall (0 .. number of walls) do                      |
	// +-------------------+------------------------------------------+
	// | 1 byte  (uint8)   | number of collider in wall               |
	// +-------------------+------------------------------------------+
	// | For each collider in wall do                                 |
	// +-------------------+------------------------------------------+
	// | 8 bytes (float64) | collider point x axis                    |
	// | 8 bytes (float64) | collider point y axis                    |
	// | 1 byte  (uint8)   | collider type                            |
	// +-------------------+------------------------------------------+
	// | End for each collider in wall                                |
	// +-------------------+------------------------------------------+
	MessageMapState = 4

	MessageGameEnd = 5
)

type MessageGameStateToEncode struct {
	CurrentTick  int32
	CurrentRound int8
	Players      []*Player
	Coins        []*Scorer
}

func (m *MessageGameStateToEncode) Encode(w codec.Writer) (err error) {
	if err = w.WriteInt32(m.CurrentTick); err != nil {
		return
	}

	if err = w.WriteInt8(m.CurrentRound); err != nil {
		return
	}

	if err = w.WriteInt32(int32(len(m.Players))); err != nil {
		return
	}

	for _, p := range m.Players {
		if err = p.Encode(w); err != nil {
			return
		}
	}

	if err = w.WriteInt32(int32(len(m.Coins))); err != nil {
		return
	}

	for _, c := range m.Coins {
		if err = c.Encode(w); err != nil {
			return
		}
	}

	return
}

type MessageGameStateToDecode struct {
	CurrentTick   int32
	CourrentRound int8
	Players       []PlayerInfo
	Coins         []struct {
		Uuid  [16]byte
		Value int32
		Pos   Point
	}
}

func (m *MessageGameStateToDecode) Decode(r codec.Reader) (err error) {
	if m.CurrentTick, err = r.ReadInt32(); err != nil {
		return
	}

	if m.CourrentRound, err = r.ReadInt8(); err != nil {
		return
	}

	var size int32
	if size, err = r.ReadInt32(); err != nil {
		return
	}

	m.Players = make([]PlayerInfo, 0, size)
	for i := int32(0); i < size; i++ {
		p := PlayerInfo{}
		if err = p.Decode(r); err != nil {
			return
		}
		m.Players = append(m.Players, p)
	}

	if size, err = r.ReadInt32(); err != nil {
		return
	}

	m.Coins = make([]struct {
		Uuid  [16]byte
		Value int32
		Pos   Point
	}, 0, size)
	for i := int32(0); i < size; i++ {
		c := struct {
			Uuid  [16]byte
			Value int32
			Pos   Point
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

		m.Coins = append(m.Coins, c)
	}

	return
}

type MessageMapStateToEncode struct {
	Map     Map
	IsAdmin bool
	Storage [100]byte
}

func (m *MessageMapStateToEncode) Encode(w codec.Writer) (err error) {
	err = m.Map.Encode(w, m.IsAdmin)
	if err != nil {
		return
	}

	_, err = w.WriteBytes(m.Storage[:])
	return
}
