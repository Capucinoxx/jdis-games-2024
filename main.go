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

func init_redis() *connector.RedisService {
  service, err := connector.NewRedisService(config.RedisAddr(), config.RedisPassword(), 1)
  if err != nil {
    panic(err)
  }

  return service
}

func main() {
	mongo := init_mongo()
  redis := init_redis()

	transport := network.NewNetwork("0.0.0.0", 8087)

	am := manager.NewAuthManager(mongo)
	nm := manager.NewNetworkManager(transport, protocol.NewBinaryProtocol())
  rm := iManager.NewRoundManager()
	gm := manager.NewGameManager(am, nm, rm, &iModel.Map{})

  rm.AddChangeStageHandler(0, &iManager.DiscoveryStage{})
  rm.AddChangeStageHandler(iManager.TicksPointRushStage, &iManager.PointRushStage{})

  sm := manager.NewScoreManager(redis, mongo)

	transport.SetRegisterFunc(gm.Register)
	transport.SetUnregisterFunc(gm.UnregisterPlayer)

	go func() {
		handler.NewHttpHandler(gm, am, sm).Handle()
		log.Fatal(gm.Init())
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigs
	log.Printf("Signal received: %s", s)
}
