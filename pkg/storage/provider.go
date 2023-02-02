package storage

import (
	"database/sql"
)

type Provider struct {
	DBConnectionString *string
}

func GetPsql[T interface{}](connection *sql.DB, tableName *string) *PsqlStorage[T] {
	return NewPsqlStorage[T](connection, tableName)
}
