package main

import (
	"fmt"
	"syscall/js"

	imodel "github.com/capucinoxx/jdis-games-2024/internal/model"

	"github.com/capucinoxx/jdis-games-2024/internal/protocol"
	"github.com/capucinoxx/jdis-games-2024/pkg/model"
)

const scale = 30.0

var players = make([]*model.Player, 0)
var proto = protocol.NewBinaryProtocol()

func arrayBufferToBytes(arrayBuffer js.Value) []byte {
	uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)
	length := uint8Array.Get("length").Int()
	bytes := make([]byte, length)
	js.CopyBytesToGo(bytes, uint8Array)
	return bytes
}

func toBytes(v js.Value) []byte {
	length := v.Length()
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = byte(v.Index(i).Int())
	}
	return b
}

func getInformations(this js.Value, args []js.Value) interface{} {
	bytes := arrayBufferToBytes(args[0])

	obj := js.Global().Get("Object").New()
	msg := proto.Decode(bytes)
	obj.Set("type", int(msg.MessageType))

	if msg.MessageType == model.MessageMapState {
		body := msg.Body.(imodel.Map)

		discreteBoard := body.DiscreteMap()
		board := js.Global().Get("Array").New()
		for i := 0; i < body.Size(); i++ {
			row := js.Global().Get("Array").New()
			for j := 0; j < body.Size(); j++ {
				row.Call("push", discreteBoard[i][j])
			}
			board.Call("push", row)
		}

		colliders := js.Global().Get("Array").New()
		for _, c := range body.Colliders() {
			collider := js.Global().Get("Array").New()
			for _, p := range c.Points {
				collider.Call("push", position(*p))
			}

			colliders.Call("push", collider)
		}

		obj.Set("map", board)
		obj.Set("walls", colliders)
	}

	if msg.MessageType == model.MessageGameState {
		body := msg.Body.(model.MessageGameStateToDecode)

		obj.Set("tick", body.CurrentTick)
		obj.Set("round", body.CourrentRound)

		players := js.Global().Get("Array").New()
		for i := 0; i < len(body.Players); i++ {
			data := body.Players[i]
			player := js.Global().Get("Object").New()
			player.Set("name", data.Nickname)
			player.Set("color", int(data.Color))
			player.Set("health", int(data.Health))
			player.Set("score", int(data.Score))
			player.Set("pos", position(data.Pos))

			if data.Dest == nil {
				player.Set("dest", position(data.Pos))
			} else {
				player.Set("dest", position(*data.Dest))
			}
			player.Set("current_weapon", int(data.CurrentWeapon))

			projectiles := js.Global().Get("Array").New()
			for _, projectile := range data.Projectiles {
				p := js.Global().Get("Object").New()
				p.Set("id", format_id(projectile.Uuid))
				p.Set("pos", position(projectile.Pos))
				p.Set("dest", position(projectile.Dest))

				projectiles.Call("push", p)
			}
			player.Set("projectiles", projectiles)
			players.Call("push", player)

			blade := js.Global().Get("Object").New()
			blade.Set("start", position(data.Blade.Start))
			blade.Set("end", position(data.Blade.End))
			blade.Set("rotation", data.Blade.Rotation)
			player.Set("blade", blade)
		}

		obj.Set("players", players)

		coins := js.Global().Get("Array").New()
		for _, coin := range body.Coins {
			c := js.Global().Get("Object").New()
			c.Set("id", format_id(coin.Uuid))
			c.Set("pos", position(coin.Pos))
			c.Set("value", coin.Value)
			coins.Call("push", c)
		}

		obj.Set("coins", coins)
	}

	return obj
}

func position(pos model.Point) interface{} {
	obj := js.Global().Get("Object").New()
	obj.Set("x", pos.X*scale)
	obj.Set("y", pos.Y*scale)
	return obj
}

func format_id(uuid [16]byte) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16],
	)
}

func registerCallbacks() {
	js.Global().Set("getInformations", js.FuncOf(getInformations))
}

func main() {
	c := make(chan struct{}, 0)

	registerCallbacks()

	<-c
}
