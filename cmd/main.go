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
	"time"
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

	srvCtx := context.WithValue(ctx, "startup_time", time.Now().UnixNano())
	problem, ok := server.NewRestServer(server.NewRestServerOptions(":3500", log)).
		Use(middleware.Otel()).
		Use(middleware.Json()).
		Router(
			server.NewRestRouter().
				Get("/", transactionModuleProvider.ProvideRoute().Find).
				POST("/", transactionModuleProvider.ProvideRoute().Create),
		).
		Start(srvCtx)

	if !ok {
		log.Fatal(context.Background(), problem.Error())
	}

}
