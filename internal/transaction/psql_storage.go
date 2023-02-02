package transaction

import (
	"database/sql"
	"github.com/yurikilian/bills/pkg/storage"
)

var tableName = "transactions"

type PsqlStorageProvider struct {
	connection *sql.DB
}

func NewPsqlStorageProvider(connection *sql.DB) *PsqlStorageProvider {
	return &PsqlStorageProvider{connection: connection}
}

func (p *PsqlStorageProvider) GetPostgres() *storage.PsqlStorage[Entity] {
	return storage.GetPsql[Entity](p.connection, &tableName)
}
