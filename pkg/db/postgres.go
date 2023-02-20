package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yurikilian/bills/internal/logger"
)

type CloseFunc func()

func ConnectPgsql(ctx context.Context, connectionString *string) (*sql.DB, CloseFunc) {
	db, err := sql.Open("pgx", *connectionString)

	if err != nil {
		logger.Log.Fatal(ctx, fmt.Sprintf("Unable to create connection pool: %v\n", err))
	}

	if err = db.Ping(); err != nil {
		logger.Log.Fatal(ctx, fmt.Sprintf("Unable to create connection pool: %v\n", err))
	}

	db.SetMaxOpenConns(10)

	return db, func() {
		err := db.Close()
		if err != nil {
			logger.Log.Error(context.Background(), err.Error())
		}
	}
}
