package transaction

import (
	"github.com/yurikilian/bills/pkg/storage"
)

type IRepository interface {
	Create(transaction *Entity) (*Entity, error)
	Find(transactionId float64) (*Entity, error)
}

type Repository struct {
	storage storage.Storage[Entity]
}

func (r *Repository) Find(transactionId float64) (*Entity, error) {
	return r.storage.Find(transactionId)
}

func (r *Repository) Create(transaction *Entity) (*Entity, error) {
	return r.storage.Create(transaction)
}

func NewRepository(storage storage.Storage[Entity]) IRepository {
	return &Repository{
		storage: storage,
	}
}
