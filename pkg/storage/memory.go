package storage

type InMemoryStorage[T any] struct {
	nextFloat float64
	memory    map[float64]*T
}

func (r *InMemoryStorage[T]) Find(id float64) (*T, error) {
	return r.memory[id], nil
}

func (r *InMemoryStorage[T]) Create(entity *T) (*T, error) {

	r.nextFloat++
	r.memory[r.nextFloat] = entity

	return entity, nil
}

func NewInMemoryStorage[T any]() *InMemoryStorage[T] {
	return &InMemoryStorage[T]{
		nextFloat: 0,
		memory:    make(map[float64]*T, 0),
	}
}

var _ Storage[interface{}] = (*InMemoryStorage[interface{}])(nil)
