package vote

import (
	"errors"
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

// CreateVote creates an vote and stores it
func (s *voteService) CreateVote(id, session, document string) (Vote, error) {
	vote := Vote{
		AssociateID: id,
		SessionID:   session,
		Document:    document,
		Creation:    time.Now(),
	}

	_, err := s.repo.FindSession(vote.SessionID)
	if err != nil {
		return Vote{}, err
	}

	err = s.repo.InsertVote(vote)
	if err != nil {
		return Vote{}, err
	}
	return vote, nil
}
