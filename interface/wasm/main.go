package main

import (
	"syscall/js"

	imodel "github.com/capucinoxx/forlorn/internal/model"

	"github.com/capucinoxx/forlorn/internal/protocol"
	"github.com/capucinoxx/forlorn/pkg/model"
)

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
	println(args[0].Length())
	bytes := arrayBufferToBytes(args[0])

	msg := proto.Decode(bytes)

	if msg.MessageType == 4 {
		body := msg.Body.(imodel.Map)

		discreteBoard := body.DiscreteMap()

		obj := js.Global().Get("Object").New()

		board := js.Global().Get("Array").New()
		for i := 0; i < body.Size(); i++ {
			row := js.Global().Get("Array").New()
			for j := 0; j < body.Size(); j++ {
				row.Call("push", discreteBoard[i][j])
			}
			board.Call("push", row)
		}

		obj.Set("Type", msg.MessageType)
		obj.Set("map", board)
		return obj
	}

	return js.ValueOf(nil)
}

func registerCallbacks() {
	js.Global().Set("getInformations", js.FuncOf(getInformations))
}

func main() {
	c := make(chan struct{}, 0)

	registerCallbacks()

	<-c
}
