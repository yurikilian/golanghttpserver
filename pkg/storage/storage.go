package storage

type Storage[T any] interface {
	Find(id float64) (*T, error)
	Create(*T) (*T, error)
}
