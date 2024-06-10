package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/capucinoxx/forlorn/internal/handler"
	iManager "github.com/capucinoxx/forlorn/internal/manager"
	iModel "github.com/capucinoxx/forlorn/internal/model"
	"github.com/capucinoxx/forlorn/internal/protocol"
	"github.com/capucinoxx/forlorn/pkg/config"
	"github.com/capucinoxx/forlorn/pkg/connector"
	"github.com/capucinoxx/forlorn/pkg/manager"
	"github.com/capucinoxx/forlorn/pkg/network"
)

func init_mongo() *connector.MongoService {
	service, err := connector.NewMongoService(config.MongoDNS(), config.MongoDatabase())
	if err != nil {
		panic(err)
	}

	return service
}

func main() {
	mongo := init_mongo()

	transport := network.NewNetwork("0.0.0.0", 8087)

	am := manager.NewAuthManager(mongo)
	nm := manager.NewNetworkManager(transport, protocol.NewBinaryProtocol())
  rm := &iManager.RoundManager{}
	gm := manager.NewGameManager(am, nm, rm, &iModel.Map{})

  // rien pour le moment
  rm.AddChangeStageHandler(0, iManager.DiscoveryStage, nil)

  // TODO: reset players avec nouvelle pos (full health, 0 bullets)
  // TODO: effacer toute les pièces
  // TODO: ajouter pièces au milieu
  rm.AddChangeStageHandler(iManager.TicksPointRushStage, iManager.PointRushStage, nil)

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
