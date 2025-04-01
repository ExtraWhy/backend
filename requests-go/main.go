package main

import (
	"casino/rest-backend/server"
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
	if _, err := conf.LoadConfig("requests-service.yaml"); err != nil {
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}
	fmt.Println("--- server up ---")
	s := server.Server{Config: &conf.Requests}
	s.DoRun()
}
