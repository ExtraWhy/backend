package main

import (
	"fmt"
	"os"
	"user-service/handlers"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/db"
)

func main() {

	yaml_config_path := "user-service.yaml"

	if len(os.Args) == 2 {
		yaml_config_path = os.Args[1]
	} else {
		fmt.Printf("No config file provided, using default: %s\n", yaml_config_path)
	}

	var conf = config.UserService{}
	if err := conf.LoadConfig(yaml_config_path); err != nil {
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	dbc := &db.DBConnection{}
	dbc.Init(conf.DBName, db.CreateUsersTable)
	dbc.Deinit()

	oauth_handler := handlers.OAuthHandler{Config: &conf}
	oauth_handler.Init(dbc)

	fmt.Println("--- user-service up ---")
}
