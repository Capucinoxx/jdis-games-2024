package protocol

import (
  "fmt"

	imodel "github.com/capucinoxx/forlorn/internal/model"
	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/model"
	p "github.com/capucinoxx/forlorn/pkg/protocol"
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
	protocol.EncodeHandlers[model.Position] = bp.encodeGameState

	protocol.DecodeHandlers[model.Spawn] = bp.decodeMapState
	protocol.DecodeHandlers[model.GameStart] = bp.decodeMapState
	protocol.DecodeHandlers[model.Position] = bp.decodeGameState
  protocol.DecodeHandlers[model.Action] = bp.decodePlayerAction

	return protocol
}

func (b BinaryProtocol) encodePlayerState(w *codec.ByteWriter, message *model.ClientMessage) {
	p := message.Body.(*model.Player)

	_ = p.Collider.Pivot.Encode(w)
	_ = w.WriteByte(byte(p.Health.Load()))
}

func (b BinaryProtocol) encodeMapState(w *codec.ByteWriter, message *model.ClientMessage) {
	p := message.Body.(*imodel.Map)

	_ = p.Encode(w)
}

func (b BinaryProtocol) encodeGameState(w *codec.ByteWriter, message *model.ClientMessage) {
  players := message.Body.([]*model.Player)

  for _, player := range players {
    fmt.Println("player encoded")
    player.Encode(w) 
  }
}

func (b BinaryProtocol) decodeMapState(r *codec.ByteReader, message *model.ClientMessage) {
	mapState := &imodel.Map{}
	if err := mapState.Decode(r); err != nil {
		message.Body = nil
	} else {
		message.Body = *mapState
	}
}

func (b BinaryProtocol) decodeGameState(r *codec.ByteReader, message *model.ClientMessage) {
  players := make([]model.PlayerInfo, 0)

  var err error
  for err == nil {
    player := model.PlayerInfo{}
    err = player.Decode(r)
    if err == nil {
      players = append(players, player)
    }
  }

  message.Body = players
}

func (b BinaryProtocol) decodePlayerAction(r *codec.ByteReader, message *model.ClientMessage) {
  var action model.Controls

  _ = r.ReadJSON(&action)

  message.Body = action
}

