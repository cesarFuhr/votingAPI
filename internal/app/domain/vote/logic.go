package vote

import (
	"errors"
	"strings"
	"time"
)

// NewVoteService creates and returns an agenda service
func NewVoteService(r Repository) Service {
	return &voteService{
		repo: r,
	}
}

type voteService struct {
	repo Repository
}

// ErrDuplicateVote represents an error caused a voting duplication
var ErrDuplicateVote = errors.New("Duplicate vote")

// ErrBadVoteFormat represents an error caused a voting duplication
var ErrBadVoteFormat = errors.New("Bad formating in vote. Must be 'S' or 'N'")

// CreateVote creates an vote and stores it
func (s *voteService) CreateVote(id, session, document, vote string) (Vote, error) {
	if len(vote) > 1 || !strings.Contains("SN", vote) {
		return Vote{}, ErrBadVoteFormat
	}

	v := Vote{
		AssociateID: id,
		SessionID:   session,
		Document:    document,
		Vote:        vote,
		Creation:    time.Now(),
	}

	_, err := s.repo.FindSession(v.SessionID)
	if err != nil {
		return Vote{}, err
	}

	err = s.repo.InsertVote(v)
	if err != nil {
		return Vote{}, err
	}
	return v, nil
}
