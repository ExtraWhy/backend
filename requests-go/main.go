package main

import (
	"casino/rest-backend/config"
	"casino/rest-backend/server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error useage : provide yaml config file")
		os.Exit(-1)
	}
	var conf = config.AppConfig{}
	if err := conf.LoadConfig(os.Args[1]); err != nil {
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}
	fmt.Println("--- server up ---")
	s := server.Server{Config: &conf}
	s.DoRun()

}
