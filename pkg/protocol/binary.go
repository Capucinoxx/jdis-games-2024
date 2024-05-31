package protocol

import (
	"encoding/binary"

	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/model"
)

// BinaryProtocol est une structure qui permet de gérer les paquets de données
// envoyés par les joueurs. Cela défini la fonction à appeler pour traiter
// l'encodage et le décodage des données selon le type de message.
type BinaryProtocol struct {
	EncodeHandlers map[model.MessageType]func(w *codec.ByteWriter, message *model.ClientMessage)
	DecodeHandlers map[model.MessageType]func(r *codec.ByteReader, message *model.ClientMessage)
}

// Encode permet d'encoder un message en un tableau d'octets.
// Le message est composé de l'identifiant du joueur, du type de message et des données à envoyer.
// représentation du message :
// [0:1 id] [1:2 messageType] [2:6 currentTime] [6:fin messageData]
func (b BinaryProtocol) Encode(id uint8, currentTime uint32, message *model.ClientMessage) []byte {
	writer := codec.NewByteWriter(binary.LittleEndian)

	// _ = writer.WriteUint8(id)
	_ = writer.WriteUint8(uint8(message.MessageType))
	// _ = writer.WriteUint32(currentTime)

	if handler, ok := b.EncodeHandlers[message.MessageType]; ok {
		handler(writer, message)
	}

	return writer.Bytes()
}

// Decode permet de décoder un tableau d'octets en un message.
// Le message est composé du tyme de message et des données reçues.
// représentation du message :
// [0:1 messageType] [1:fin messageData]
func (b BinaryProtocol) Decode(data []byte) model.ClientMessage {
	reader := codec.NewByteReader(data[1:], binary.LittleEndian)

	msg := model.ClientMessage{
		MessageType: model.MessageType(data[0]),
	}

	if handler, ok := b.DecodeHandlers[msg.MessageType]; ok {
		// handler(data[1:], &msg)
		handler(reader, &msg)
	}

	return msg
}
