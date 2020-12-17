package session

import (
	"time"

	"github.com/google/uuid"
)

// NewSessionService creates and returns an agenda service
func NewSessionService(r Repository, p Publisher) Service {
	return &sessionService{
		repo:  r,
		pub:   p,
		clock: &internalClock{},
	}
}

type sessionService struct {
	repo  Repository
	pub   Publisher
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

	t := time.NewTimer(session.Duration + (10 * time.Second))
	go notifyResult(t, session, s)

	return session, nil
}

func notifyResult(t *time.Timer, session Session, s *sessionService) {
	<-t.C
	defer t.Stop()

	result, err := s.Result(session.ID)
	if err != nil {
		return
	}
	err = s.pub.PublishResult(result)
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

	c := Count{}
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
