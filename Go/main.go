package main

import (
	"backend/config"
	"backend/utils/helper"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

func init() {
	config.InitEnv("./.env")
	config.ConfigureLogger()
}

var logger = logging.MustGetLogger("main")

// @title Widatech Test
func main() {
	logger.Debugf("Envs: %v", helper.PrettyJson(config.Envs))

	gormDB := config.NewPostgresqlDB()

	// migrations
	err := gormDB.AutoMigrate()
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// repositories

	// ucases

	dependencies := CommonDeps{}

	ginEngine := gin.Default()
	SetupServer(ginEngine, dependencies)
	ginEngine.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
}
