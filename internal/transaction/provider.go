package transaction

import (
	"database/sql"
	"github.com/yurikilian/bills/pkg/storage"
)

type ModuleProvider struct {
	storage storage.Storage[Entity]
	route   *Route
}

func (p *ModuleProvider) ProvideRoute() *Route {
	if p.route != nil {
		return p.route
	}

	p.route = &Route{service: NewTransactionService(NewRepository(p.storage))}
	return p.route
}

type ModuleBuilder struct {
	provider *ModuleProvider
}

func NewTransactionModuleBuilder() *ModuleBuilder {
	return &ModuleBuilder{
		provider: &ModuleProvider{},
	}
}

func (p *ModuleBuilder) WithPsqlStorage(dbConnection *sql.DB) *ModuleBuilder {

	if p.provider.storage != nil {
		panic("Storage already defined")
	}

	tableName := "transaction"
	p.provider.storage = storage.GetPsql[Entity](dbConnection, &tableName)
	return p
}

func (p *ModuleBuilder) WithInMemoryStorage(inMemory *storage.InMemoryStorage[Entity]) *ModuleBuilder {
	if p.provider.storage != nil {
		panic("Storage already defined")
	}
	p.provider.storage = inMemory
	return p
}

func (p *ModuleBuilder) Build() *ModuleProvider {
	if p.provider.storage == nil {
		panic("the transaction storage is not defined")
	}

	return p.provider
}
