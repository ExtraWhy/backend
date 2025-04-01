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

	var conf = config.MegaConfig{}
	service_config := config.UserService{}

	if err := conf.LoadConfig(yaml_config_path, &service_config); err != nil {
		fmt.Println("Failed to load cofig file")
		os.Exit(-2)
	}

	dbc := &db.DBConnection{}
	dbc.Init(service_config.DBName, service_config.DBDriver)
	if err := dbc.Init(service_config.DBDriver, service_config.DBName); err != nil {
		fmt.Printf("Failed to initialize DB: %v", err)
	}
	dbc.SetupSchema(db.CreateUsersTable)
	defer dbc.Deinit()

	oauth_handler := handlers.OAuthHandler{Config: &service_config}
	oauth_handler.Init(dbc)

	fmt.Println("--- user-service up ---")
}
