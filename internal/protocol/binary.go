package protocol

import (
	"encoding/binary"
	"math"

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
		EncodeHandlers: make(map[model.MessageType]func(message *model.ClientMessage) []byte),
		DecodeHandlers: make(map[model.MessageType]func(data []byte, message *model.ClientMessage)),
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
func (b BinaryProtocol) encodePlayerState(message *model.ClientMessage) []byte {
	p := message.Body.(*model.Player)
	buf := make([]byte, PlayerPacketSize)

	binary.LittleEndian.PutUint32(buf[0:4], math.Float32bits(p.Collider.Pivot.X))
	binary.LittleEndian.PutUint32(buf[4:8], math.Float32bits(p.Collider.Pivot.Y))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(p.Collider.Rotation))

	// TODO: shooting statement

	buf[13] = byte(p.Health)

	return buf
}

// encodeMapState permet d'encoder l'état de la map.
// L'état de la map est composé de tous les colliders présents
// dans la map.
func (b BinaryProtocol) encodeMapState(message *model.ClientMessage) []byte {
	p := message.Body.([]*model.Collider)
	buf := make([]byte, 0)

	for _, c := range p {
		buf = append(buf, b.encodeCollider(c)...)
	}

	return buf
}

// encodeCollider permet d'encoder un collider.
// Voici la représentation d'un collider :
// [0:1 type] [1:2 taille] [[p:p+4 x] [p+4:p+8 y] ...]
func (b BinaryProtocol) encodeCollider(c *model.Collider) []byte {
	buf := make([]byte, 0)

	buf = append(buf, byte(c.Type))
	buf = append(buf, byte(len(c.Points)))

	for _, p := range c.Points {
		binary.LittleEndian.PutUint32(buf, uint32(p.X))
		binary.LittleEndian.PutUint32(buf, uint32(p.Y))
	}

	return buf
}

func decodePoint(data []byte) *model.Point {
	p := &model.Point{
		X: math.Float32frombits(binary.LittleEndian.Uint32(data[4:])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(data[4:8])),
	}

	if math.IsNaN(float64(p.X)) || math.IsNaN(float64(p.Y)) {
		return nil
	}

	return p
}

func decodePlayerInput(data []byte, message *model.ClientMessage) {
	controls := model.Controls{
		Rotation: binary.LittleEndian.Uint32(data),
		Shoot:    decodePoint(data[4:12]),
	}

	// decode shoot position
	message.Body = controls
}
