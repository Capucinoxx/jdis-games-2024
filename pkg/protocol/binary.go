package protocol

import (
	"encoding/binary"

	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/model"
)

type BinaryProtocol struct {
	EncodeHandlers map[model.MessageType]func(w *codec.ByteWriter, message *model.ClientMessage)
	DecodeHandlers map[model.MessageType]func(r *codec.ByteReader, message *model.ClientMessage)
}

func (b BinaryProtocol) Encode(id uint8, currentTime uint32, message *model.ClientMessage) []byte {
	writer := codec.NewByteWriter(binary.LittleEndian)

	_ = writer.WriteUint8(uint8(message.MessageType))

	if handler, ok := b.EncodeHandlers[message.MessageType]; ok {
		handler(writer, message)
	}

	return writer.Bytes()
}

func (b BinaryProtocol) Decode(data []byte) model.ClientMessage {
	reader := codec.NewByteReader(data[1:], binary.LittleEndian)

	msg := model.ClientMessage{
		MessageType: model.MessageType(data[0]),
	}

	if handler, ok := b.DecodeHandlers[msg.MessageType]; ok {
		handler(reader, &msg)
	}

	return msg
}
