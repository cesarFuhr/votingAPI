package vote

import (
	"errors"
	"strings"
	"time"
)

// NewVoteService creates and returns an agenda service
func NewVoteService(r Repository, v DocValidator) Service {
	return &voteService{
		repo:      r,
		validator: v,
		clock:     &internalClock{},
	}
}

type internalClock struct{}

func (c *internalClock) Now() time.Time {
	return time.Now()
}

type clock interface {
	Now() time.Time
}

type voteService struct {
	repo      Repository
	validator DocValidator
	clock     clock
}

var (
	// ErrDuplicateVote represents an error caused a voting duplication
	ErrDuplicateVote = errors.New("Duplicate vote")
	// ErrBadVoteFormat represents an error caused a voting duplication
	ErrBadVoteFormat = errors.New("Bad formating in vote. Must be 'S' or 'N'")
	// ErrSessionExpired represents an error caused by session expiration
	ErrSessionExpired = errors.New("This voting session is expired")
	// ErrNotAbleToVote represents an error caused by invalid document
	ErrNotAbleToVote = errors.New("Associate not able to vote")
)

// CreateVote creates an vote and stores it
func (s *voteService) CreateVote(id, session, document, vote string) (Vote, error) {
	if len(vote) > 1 || !strings.Contains("SN", vote) {
		return Vote{}, ErrBadVoteFormat
	}

	isValidDoc, err := s.validator.ValidateDocument(document)
	if err != nil {
		return Vote{}, ErrSessionExpired
	}
	if !isValidDoc {
		return Vote{}, ErrNotAbleToVote
	}

	v := Vote{
		AssociateID: id,
		SessionID:   session,
		Document:    document,
		Vote:        vote,
		Creation:    time.Now(),
	}

	sess, err := s.repo.FindSession(v.SessionID)
	if err != nil {
		return Vote{}, err
	}

	if s.clock.Now().After(sess.GetExpiration()) {
		return Vote{}, ErrSessionExpired
	}

	err = s.repo.InsertVote(v)
	if err != nil {
		return Vote{}, err
	}
	return v, nil
}
