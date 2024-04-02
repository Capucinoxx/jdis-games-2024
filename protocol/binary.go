package protocol

import (
	"encoding/binary"

	"github.com/capucinoxx/forlorn/model"
)

const (
	// PlayerPacketSize représente la taille en octets
	// nécessaire pour stocker un paquet de données d'un joueur
	PlayerPacketSize = 14
)

// BinaryProtocol est une structure qui permet de gérer les paquets de données
// envoyés par les joueurs. Cela défini la fonction à appeler pour traiter
// l'encodage et le décodage des données selon le type de message.
type BinaryProtocol struct {
	encodeHandlers map[model.MessageType]func(message *model.ClientMessage) []byte
	decodeHandlers map[model.MessageType]func(data []byte, message *model.ClientMessage)
}

// NewBinaryProtocol crée un nouveau protocole binaire. Ce protocole
// permet de gérer les messages clients en les encodant et les décodant
// en un tableau d'octets.
func NewBinaryProtocol() *BinaryProtocol {
	protocol := &BinaryProtocol{
		encodeHandlers: make(map[model.MessageType]func(message *model.ClientMessage) []byte),
		decodeHandlers: make(map[model.MessageType]func(data []byte, message *model.ClientMessage)),
	}

	protocol.encodeHandlers[0] = protocol.encodePlayerState

	protocol.decodeHandlers[0] = decodePlayerTime
	protocol.decodeHandlers[1] = decodePlayerInput

	return protocol
}

// Encode permet d'encoder un message en un tableau d'octets.
// Le message est composé de l'identifiant du joueur, du type de message et des données à envoyer.
// représentation du message :
// [0:1 id] [1:2 messageType] [2:6 currentTime] [6:fin messageData]
func (b BinaryProtocol) Encode(id uint8, currentTime uint32, message *model.ClientMessage) []byte {
	buf := make([]byte, 0)
	buf = append(buf, byte(id))
	buf = append(buf, byte(message.MessageType))

	currentTimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(currentTimeBytes, currentTime)
	buf = append(buf, currentTimeBytes...)

	if handler, ok := b.encodeHandlers[message.MessageType]; ok {
		buf = append(buf, handler(message)...)
	}

	return buf
}

// encodePlayerState permet d'encoder l'état d'un joueur.
// L'état d'un joueur est composé de sa position, de son orientation,
// de son état de tir et de sa vie.
// représentation de l'état d'un joueur :
// [0:4 x] [4:8 y] [8:12 rotation] [12:13 shooting] [13:14 health]
func (b BinaryProtocol) encodePlayerState(message *model.ClientMessage) []byte {
	p := message.Body.(*model.Player)
	buf := make([]byte, PlayerPacketSize)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(p.Collider.Pivot.X))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(p.Collider.Pivot.Y))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(p.Collider.Rotation))

	// TODO: shooting statement

	buf[13] = byte(p.Health)

	return buf
}

// Decode permet de décoder un tableau d'octets en un message.
// Le message est composé du tyme de message et des données reçues.
// représentation du message :
// [0:1 messageType] [1:fin messageData]
func (b BinaryProtocol) Decode(data []byte) model.ClientMessage {
	msg := model.ClientMessage{
		MessageType: model.MessageType(data[0]),
	}

	if handler, ok := b.decodeHandlers[msg.MessageType]; ok {
		handler(data[1:], &msg)
	}

	return msg
}

// decodePlayerTime permet de décoder le temps envoyé par un joueur.
func decodePlayerTime(data []byte, message *model.ClientMessage) {
	message.Body = binary.LittleEndian.Uint32(data)
}

// decodePlayerInput permet de décoder les données de contrôle
// envoyées par un joueur.
func decodePlayerInput(data []byte, message *model.ClientMessage) {
	p := message.Body.(*model.Player)

	p.Controls.Left = data[0] == 1
	p.Controls.Right = data[1] == 1

	message.Body = p
}
