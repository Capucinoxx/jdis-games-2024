package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/capucinoxx/forlorn/internal/handler"
	"github.com/capucinoxx/forlorn/internal/protocol"
	"github.com/capucinoxx/forlorn/pkg/manager"
	"github.com/capucinoxx/forlorn/pkg/network"
)

func main() {
	transport := network.NewNetwork("localhost", 8087)

	am := manager.NewAuthManager()
	nm := manager.NewNetworkManager(transport, protocol.NewBinaryProtocol())
	gm := manager.NewGameManager(am, nm)

	transport.SetRegisterFunc(gm.RegisterPlayer)
	transport.SetUnregisterFunc(gm.UnregisterPlayer)

	go func() {
		handler.NewHttpHandler(gm, am).Handle()
		log.Fatal(gm.Init())
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigs
	log.Printf("Signal received: %s", s)
}
