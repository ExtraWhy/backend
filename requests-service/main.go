package main

import (
	server "casino/rest-server"
	servinterface "casino/serv-interface"
	websocket "casino/ws-server"
	"fmt"
	"os"
	"sync"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/logger"
	"go.uber.org/zap"
)

var (
	version string = "local" // automatically populated by the build system

	zl = logger.ZapperLog{}
	do sync.Once
)

func log(level int, m string, zpf ...zap.Field) {
	do.Do(func() {
		zl.Init(logger.DEV)
	})

	zpf = append(zpf, zap.String("version", version))
	zl.Log(level, m, zpf...)
}

func main() {
	if len(os.Args) != 2 {
		log(logger.CRITICAL, "Error usage: provide config file")
		//	os.Exit(-1)
	}

	var conf = config.MegaConfig{}
	srvconf := config.RequestService{}
	//if err := conf.LoadConfig("requests-service.yaml", &srvconf); err != nil {
	if err := conf.LoadConfig(os.Args[1], &srvconf); err != nil {
		log(logger.CRITICAL, "Failed to load config file", zap.Any("what", err))
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		var srvIface servinterface.ServiceInterface
		log(logger.INFO, "--- rest service up ---")
		srvIface = &server.Server{}
		if err := srvIface.DoRun(&srvconf); err != nil {
			log(logger.CRITICAL, "Failed run rest service", zap.Any("what", err))
			os.Exit(-3)
		}
	}()

	go func() {
		defer wg.Done()
		var srvIface servinterface.ServiceInterface

		log(logger.INFO, "--- websocket service up ---")

		srvIface = &websocket.WSServer{}

		if err := srvIface.DoRun(&srvconf); err != nil {
			log(logger.CRITICAL, "Failed run websocket service", zap.Any("what", err))
			os.Exit(-1)
		}
	}()

	wg.Wait()
}
