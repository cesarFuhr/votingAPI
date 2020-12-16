package agenda

import "github.com/google/uuid"

// NewAgendaService creates and returns an agenda service
func NewAgendaService(r Repository) Service {
	return &agendaService{
		repo: r,
	}
}

type agendaService struct {
	repo Repository
}

// CreateAgenda creates an agenda em stores it
func (s *agendaService) CreateAgenda(description string) (Agenda, error) {
	id := uuid.New()

	agenda := Agenda{
		ID:          id.String(),
		Description: description,
	}

	err := s.repo.InsertAgenda(agenda)
	if err != nil {
		return Agenda{}, err
	}
	return agenda, nil
}

// FindAgenda returns a agenda finding by ID
func (s *agendaService) FindAgenda(id string) (Agenda, error) {
	agenda, err := s.repo.FindAgenda(id)
	if err != nil {
		return Agenda{}, err
	}
	return agenda, nil
}
