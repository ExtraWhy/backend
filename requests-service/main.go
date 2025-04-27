package main

import (
	server "casino/rest-backend/rest-server"
	servinterface "casino/rest-backend/serv-interface"
	websocket "casino/rest-backend/ws-server"
	"fmt"
	"os"
	"sync"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/logger"
	"go.uber.org/zap"
)

var (
	zl = logger.ZapperLog{}
	do sync.Once
)

func log(level int, m string, zpf ...zap.Field) {
	do.Do(func() {
		zl.Init(logger.DEV)
	})
	zl.Log(level, m, zpf...)
}

/* //DEBUG
func main() {

	var srvIface servinterface.ServiceInterface

	var conf = config.MegaConfig{}
	srvconf := config.RequestService{}
	//conf.LoadConfig("requests-service.yaml", &req); err != nil {
	if err := conf.LoadConfig("requests-service.yaml", &srvconf); err != nil {
		log(logger.CRITICAL, "Failed to load config file", zap.Any("what", err))
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	if srvconf.ApiType == "rest" {
		log(logger.INFO, "--- rest service up ---")
		srvIface = &server.Server{}
		srvIface.DoRun(&srvconf)
	} else if srvconf.ApiType == "ws" {
		log(logger.INFO, "--- websocket service up ---")
		srvIface = &websocket.WSServer{}
		srvIface.DoRun(&srvconf)
	}
}
*/

// MAIN
func main() {

	if len(os.Args) != 2 {
		log(logger.CRITICAL, "Error usage: provide config file")
		os.Exit(-1)
	}

	var srvIface servinterface.ServiceInterface

	var conf = config.MegaConfig{}
	srvconf := config.RequestService{}
	//conf.LoadConfig("requests-service.yaml", &req); err != nil {
	if err := conf.LoadConfig(os.Args[1], &srvconf); err != nil {
		log(logger.CRITICAL, "Failed to load config file", zap.Any("what", err))
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	if srvconf.ApiType == "rest" {
		log(logger.INFO, "--- rest service up ---")
		srvIface = &server.Server{}
		srvIface.DoRun(&srvconf)
	} else if srvconf.ApiType == "ws" {
		log(logger.INFO, "--- websocket service up ---")
		srvIface = &websocket.WSServer{}
		srvIface.DoRun(&srvconf)
	}
}
