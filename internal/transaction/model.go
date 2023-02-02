package transaction

type CreationRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,number"`
	Currency    string  `json:"currency" validate:"required,oneof=EUR"`
	Type        string  `json:"type" validate:"required,oneof=CREDIT DEBIT"`
}

type FindRequest struct {
	id float64 `param:"id"`
}
