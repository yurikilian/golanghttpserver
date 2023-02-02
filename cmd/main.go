package main

import (
	"context"
	_ "embed"
	"github.com/yurikilian/bills/internal/logger"
	"github.com/yurikilian/bills/internal/transaction"
	"github.com/yurikilian/bills/pkg/db"
	"github.com/yurikilian/bills/pkg/middleware"
	"github.com/yurikilian/bills/pkg/server"
	"github.com/yurikilian/bills/pkg/telemetry"
)

func main() {

	ctx := context.Background()
	closeTelemetry := telemetry.Init(ctx)
	defer closeTelemetry()

	log := logger.NewProvider().ProvideLog()

	configurationProvider := server.NewConfigurationProvider()
	dbConnection := db.ConnectPgsql(configurationProvider.GetDBConnectionString())

	transactionModuleProvider := transaction.
		NewTransactionModuleBuilder().
		WithPsqlStorage(dbConnection).
		Build()

	err := server.NewRestServer(server.NewRestServerOptions(":8080", log)).
		Use(middleware.Otel()).
		Use(middleware.Json()).
		Router(
			server.NewRestRouter().
				Get("/", transactionModuleProvider.ProvideRoute().Find).
				POST("/", transactionModuleProvider.ProvideRoute().Create),
		).
		Start()

	if err != nil {
		log.Fatal(context.Background(), err.Error())
	}

}
