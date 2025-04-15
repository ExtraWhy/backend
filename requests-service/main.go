package main

import (
	server "casino/rest-backend/rest-server"

	"fmt"
	"os"

	"github.com/ExtraWhy/internal-libs/config"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error useage : provide yaml config file")
		os.Exit(-1)
	}

	var conf = config.MegaConfig{}
	srvconf := config.RequestService{}
	//conf.LoadConfig("requests-service.yaml", &req); err != nil {
	if err := conf.LoadConfig(os.Args[1], &srvconf); err != nil {
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	fmt.Println("--- server up ---")
	s := server.Server{}
	s.DoRun(&srvconf)
}
