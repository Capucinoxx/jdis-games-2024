package protocol

import (
	imodel "github.com/capucinoxx/jdis-games-2024/internal/model"
	"github.com/capucinoxx/jdis-games-2024/pkg/codec"
	"github.com/capucinoxx/jdis-games-2024/pkg/model"
	p "github.com/capucinoxx/jdis-games-2024/pkg/protocol"
)

type BinaryProtocol struct{}

func NewBinaryProtocol() *p.BinaryProtocol {
	protocol := &p.BinaryProtocol{
		EncodeHandlers: make(map[model.MessageType]func(w *codec.ByteWriter, message *model.ClientMessage)),
		DecodeHandlers: make(map[model.MessageType]func(r *codec.ByteReader, message *model.ClientMessage)),
	}

	bp := BinaryProtocol{}
	protocol.EncodeHandlers[model.MessageMapState] = bp.encodeMapState
	protocol.EncodeHandlers[model.MessageGameEnd] = bp.encodeGameEnd
	protocol.EncodeHandlers[model.MessageGameState] = bp.encodeGameState

	protocol.DecodeHandlers[model.MessageMapState] = bp.decodeMapState
	protocol.DecodeHandlers[model.MessageGameEnd] = bp.decodeGameEnd
	protocol.DecodeHandlers[model.MessageGameState] = bp.decodeGameState
	protocol.DecodeHandlers[model.MessagePlayerAction] = bp.decodePlayerAction

	return protocol
}

func (b BinaryProtocol) encodeMapState(w *codec.ByteWriter, message *model.ClientMessage) {
	p := message.Body.(model.MessageMapStateToEncode)

	_ = p.Encode(w)
}

func (b BinaryProtocol) encodeGameState(w *codec.ByteWriter, message *model.ClientMessage) {
	data := message.Body.(model.MessageGameStateToEncode)

	_ = data.Encode(w)
}

func (b BinaryProtocol) encodeGameEnd(w *codec.ByteWriter, message *model.ClientMessage) {}

func (b BinaryProtocol) decodeGameEnd(r *codec.ByteReader, message *model.ClientMessage) {}

func (b BinaryProtocol) decodeMapState(r *codec.ByteReader, message *model.ClientMessage) {
	mapState := &imodel.Map{}
	if err := mapState.Decode(r); err != nil {
		message.Body = nil
	} else {
		message.Body = *mapState
	}
}

func (b BinaryProtocol) decodeGameState(r *codec.ByteReader, message *model.ClientMessage) {
	var state model.MessageGameStateToDecode
	state.Decode(r)

	message.Body = state
}

func (b BinaryProtocol) decodePlayerAction(r *codec.ByteReader, message *model.ClientMessage) {
	var action model.Controls

	_ = r.ReadJSON(&action)

	message.Body = action
}
