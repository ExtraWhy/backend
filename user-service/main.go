package main

import (
	"fmt"
	"os"
	"sync"
	"user-service/handlers"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/db"

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
		zl.Init(1)
	})

	zpf = append(zpf, zap.String("version", version))
	zl.Log(level, m, zpf...)
}

func main() {
	log(logger.INFO, "Starting user-service")

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

	dbc := &db.DBSqlConnection{}
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
