package session

import (
	"time"

	"github.com/google/uuid"
)

// NewSessionService creates and returns an agenda service
func NewSessionService(r Repository) Service {
	return &sessionService{
		repo:  r,
		clock: &internalClock{},
	}
}

type sessionService struct {
	repo  Repository
	clock clock
}

type internalClock struct{}

func (c *internalClock) Now() time.Time {
	return time.Now()
}

type clock interface {
	Now() time.Time
}

// CreateSession creates an agenda em stores it
func (s *sessionService) CreateSession(agendaID string, duration time.Duration) (Session, error) {
	id := uuid.New()

	session := Session{
		ID:             id.String(),
		OriginalAgenda: agendaID,
		Duration:       duration,
		Creation:       s.clock.Now(),
	}

	err := s.repo.InsertSession(session)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

// FindSession returns a agenda finding by ID
func (s *sessionService) FindSession(id string) (Session, error) {
	agenda, err := s.repo.FindSession(id)
	if err != nil {
		return Session{}, err
	}
	return agenda, nil
}
