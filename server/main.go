package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/capucinoxx/jdis-games-2024/consts"
	"github.com/capucinoxx/jdis-games-2024/internal/handler"
	iManager "github.com/capucinoxx/jdis-games-2024/internal/manager"
	iModel "github.com/capucinoxx/jdis-games-2024/internal/model"
	"github.com/capucinoxx/jdis-games-2024/internal/protocol"
	"github.com/capucinoxx/jdis-games-2024/pkg/config"
	"github.com/capucinoxx/jdis-games-2024/pkg/connector"
	"github.com/capucinoxx/jdis-games-2024/pkg/manager"
	"github.com/capucinoxx/jdis-games-2024/pkg/network"
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

	transport := network.NewNetwork("0.0.0.0", config.Port())

	sm := manager.NewScoreManager(redis, mongo)

	am := manager.NewAuthManager(mongo)
	am.SetupAdmins(config.RequiredAdmins())

	nm := manager.NewNetworkManager(transport, protocol.NewBinaryProtocol())
	rm := iManager.NewRoundManager()
	gm := manager.NewGameManager(am, nm, rm, sm, &iModel.Map{})

	rm.AddChangeStageHandler(0, &iManager.DiscoveryStage{})
	rm.AddChangeStageHandler(consts.TicksPointRushStage, &iManager.PointRushStage{})

	transport.SetRegisterFunc(gm.RegisterConnection)
	transport.SetUnregisterFunc(gm.RemoveConnection)

	go func() {
		handler.NewHttpHandler(gm, am, sm).Handle()
		log.Fatal(gm.Initialize())
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigs
	log.Printf("Signal received: %s", s)
}
