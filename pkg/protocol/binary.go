package protocol

import (
	"encoding/binary"

	"github.com/capucinoxx/forlorn/pkg/model"
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
	EncodeHandlers map[model.MessageType]func(message *model.ClientMessage) []byte
	DecodeHandlers map[model.MessageType]func(data []byte, message *model.ClientMessage)
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

	if handler, ok := b.EncodeHandlers[message.MessageType]; ok {
		buf = append(buf, handler(message)...)
	}

	return buf
}

// Decode permet de décoder un tableau d'octets en un message.
// Le message est composé du tyme de message et des données reçues.
// représentation du message :
// [0:1 messageType] [1:fin messageData]
func (b BinaryProtocol) Decode(data []byte) model.ClientMessage {
	// TODO: validate data length
	msg := model.ClientMessage{
		MessageType: model.MessageType(data[0]),
	}

	if handler, ok := b.DecodeHandlers[msg.MessageType]; ok {
		handler(data[1:], &msg)
	}

	return msg
}
