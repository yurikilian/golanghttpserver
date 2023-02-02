package transaction

type EntityId[T any] interface {
	Get() *T
}

type FloatEntityId struct {
	EntityId[float64]
	id float64
}

func (f FloatEntityId) Get() *float64 {
	return &f.id
}

type Entity struct {
	Id          float64 `column:"id"`
	Title       string  `column:"title"`
	Description string  `column:"description"`
	Price       float64 `column:"price"`
	Currency    string  `column:"currency"`
	Type        string  `column:"type"`
}
