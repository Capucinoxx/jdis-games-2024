package main

import (
	"syscall/js"

	"github.com/capucinoxx/forlorn/internal/protocol"
	"github.com/capucinoxx/forlorn/pkg/model"
)

var players = make([]*model.Player, 0)

func getInformations(this js.Value, args []js.Value) interface{} {
	println(args)

	return js.ValueOf(nil)
}

func registerCallbacks() {
	js.Global().Set("getInformations", js.FuncOf(getInformations))
}

var proto protocol.BinaryProtocol

func main() {
	c := make(chan struct{}, 0)

	registerCallbacks()

	<-c
}
