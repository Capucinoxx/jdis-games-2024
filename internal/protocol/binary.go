package protocol

import (
	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/model"
	p "github.com/capucinoxx/forlorn/pkg/protocol"
)

const (
	// PlayerPacketSize représente la taille en octets
	// nécessaire pour stocker un paquet de données d'un joueur
	PlayerPacketSize = 14
)

// BinaryProtocol est une structure vide encapsulant les différentes
// fonctions de traitement des messages.
type BinaryProtocol struct{}

// NewBinaryProtocol crée un nouveau protocole binaire. Ce protocole
// permet de gérer les messages clients en les encodant et les décodant
// en un tableau d'octets.
func NewBinaryProtocol() *p.BinaryProtocol {
	protocol := &p.BinaryProtocol{
		EncodeHandlers: make(map[model.MessageType]func(w *codec.ByteWriter, message *model.ClientMessage)),
		DecodeHandlers: make(map[model.MessageType]func(r *codec.ByteReader, message *model.ClientMessage)),
	}

	bp := BinaryProtocol{}
	protocol.EncodeHandlers[model.Spawn] = bp.encodeMapState
	protocol.EncodeHandlers[model.GameStart] = bp.encodeMapState
	protocol.EncodeHandlers[model.Position] = bp.encodePlayerState

	protocol.DecodeHandlers[model.Position] = decodePlayerInput

	return protocol
}

// encodePlayerState permet d'encoder l'état d'un joueur.
// L'état d'un joueur est composé de sa position, de son orientation,
// de son état de tir et de sa vie.
// représentation de l'état d'un joueur :
// [0:4 x] [4:8 y] [8:12 rotation] [12:13 shooting] [13:14 health]
func (b BinaryProtocol) encodePlayerState(w *codec.ByteWriter, message *model.ClientMessage) {
	p := message.Body.(*model.Player)

	_ = p.Collider.Pivot.Encode(w)
	_ = w.WriteByte(byte(p.Health.Load()))
}

// encodeMapState permet d'encoder l'état de la map.
// L'état de la map est composé de tous les colliders présents
// dans la map.
func (b BinaryProtocol) encodeMapState(w *codec.ByteWriter, message *model.ClientMessage) {
	p := message.Body.([]*model.Collider)

	for _, c := range p {
		c.Encode(w)
	}
}

func decodePlayerInput(r *codec.ByteReader, message *model.ClientMessage) {
	shootPos := &model.Point{}
	if err := shootPos.Decode(r); err != nil {
		shootPos = nil
	}

	controls := model.Controls{
		Shoot: shootPos,
	}

	message.Body = controls
}
