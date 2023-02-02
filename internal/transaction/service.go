package transaction

type Service struct {
	repository IRepository
}

func (s *Service) Create(req CreationRequest) error {
	entity := &Entity{
		Title:       req.Title,
		Description: req.Description,
		Currency:    req.Currency,
		Type:        req.Type,
		Price:       req.Price,
	}
	if _, err := s.repository.Create(entity); err != nil {
		return err
	}
	return nil
}

func (s *Service) Find(id float64) (*Entity, error) {
	trn, err := s.repository.Find(id)
	if err != nil {
		return nil, err
	}
	return trn, nil
}

func NewTransactionService(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}
