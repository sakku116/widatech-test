package main

import (
	"backend/config"
	"backend/domain/model"
	"backend/repository"
	ucase "backend/usecase"
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
	err := gormDB.AutoMigrate(
		&model.Invoice{},
		&model.Product{},
	)
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// repositories
	invoiceRepo := repository.NewInvoiceRepo(gormDB)
	productRepo := repository.NewProductRepo(gormDB)

	// ucases
	invoiceUcase := ucase.NewInvoiceUcase(invoiceRepo, productRepo)

	dependencies := CommonDeps{
		invoiceUcase: invoiceUcase,
	}

	ginEngine := gin.Default()
	SetupServer(ginEngine, dependencies)
	ginEngine.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
}
