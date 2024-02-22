package model

import "time"

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

type Player struct{}
