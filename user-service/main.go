package main

import (
	"fmt"
	"os"
	"user-service/handlers"

	"github.com/ExtraWhy/internal-libs/config"
)

func main() {

	yaml_config_path := "user-service.yaml"

	if len(os.Args) == 2 {
		yaml_config_path = os.Args[1]
	} else {
		fmt.Printf("No config file provided, using default: %s\n", yaml_config_path)
	}

	var conf = config.MegaConfig{}
	usr := config.UserService{}
	if err := conf.LoadConfig(yaml_config_path, &usr); err != nil {
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	oauth_handler := handlers.OAuthHandler{Config: &usr}
	oauth_handler.Init()

	fmt.Println("--- user-service up ---")
}
