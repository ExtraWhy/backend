package main

import (
	server "casino/rest-backend/rest-server"
	servinterface "casino/rest-backend/serv-interface"
	websocket "casino/rest-backend/ws-server"

	"fmt"
	"os"

	"github.com/ExtraWhy/internal-libs/config"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error useage : provide yaml config file")
		os.Exit(-1)
	}

	var srvIface servinterface.ServiceInterface

	var conf = config.MegaConfig{}
	srvconf := config.RequestService{}
	//conf.LoadConfig("requests-service.yaml", &req); err != nil {
	if err := conf.LoadConfig(os.Args[1], &srvconf); err != nil {
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	if srvconf.ApiType == "rest" {
		fmt.Println("--- rest server up ---")
		srvIface = &server.Server{}
		srvIface.DoRun(&srvconf)
	} else if srvconf.ApiType == "ws" {
		fmt.Println("--- ws server up ---")
		srvIface = &websocket.WSServer{}
		srvIface.DoRun(&srvconf)
	}
}
