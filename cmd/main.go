package main

import (
	"context"
	_ "embed"
	"github.com/yurikilian/bills/internal/logger"
	"github.com/yurikilian/bills/internal/transaction"
	"github.com/yurikilian/bills/pkg/db"
	"github.com/yurikilian/bills/pkg/middleware"
	"github.com/yurikilian/bills/pkg/server"
	"time"
)

func main() {

	ctx := context.Background()

	configurationProvider := server.NewConfigurationProvider()
	dbConnection, closeDb := db.ConnectPgsql(ctx, configurationProvider.GetDBConnectionString())
	defer closeDb()

	/*closeTelemetry := telemetry.Init(ctx)
	defer closeTelemetry()*/

	transactionModuleProvider := transaction.
		NewTransactionModuleBuilder().
		WithPsqlStorage(dbConnection).
		Build()

	srvCtx := context.WithValue(ctx, "startup_time", time.Now().UnixNano())
	problem, ok := server.NewRestServer(server.NewRestServerOptions(":3500", logger.Log)).
		Use(middleware.Otel()).
		Use(middleware.Json()).
		Router(
			server.NewRestRouter().
				Get("/", transactionModuleProvider.ProvideRoute().Find).
				POST("/", transactionModuleProvider.ProvideRoute().Create),
		).
		Start(srvCtx)

	if !ok {
		logger.Log.Fatal(context.Background(), problem.Error())
	}

}
