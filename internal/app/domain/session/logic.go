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

// CreateSession creates an session em stores it
func (s *sessionService) CreateSession(agendaID string, duration time.Duration) (Session, error) {
	id := uuid.New()

	if duration == 0 {
		duration = time.Minute
	}

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

// FindSession returns a session finding by ID
func (s *sessionService) FindSession(id string) (Session, error) {
	session, err := s.repo.FindSession(id)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

// Result returns a voting session result
func (s *sessionService) Result(id string) (Result, error) {
	session, err := s.repo.FindSession(id)
	if err != nil {
		return Result{}, err
	}

	votes, err := s.repo.FindVotes(session)
	if err != nil {
		return Result{}, err
	}

	c := count{}
	for _, v := range votes {
		if v == "S" {
			c.InFavor++
		} else {
			c.Against++
		}
	}

	return Result{
		ID:             session.ID,
		OriginalAgenda: session.OriginalAgenda,
		Closed:         s.clock.Now().After(session.GetExpiration()),
		Count:          c,
	}, nil
}
