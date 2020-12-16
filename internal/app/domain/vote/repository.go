package vote

import "github.com/cesarFuhr/gocrypto/internal/app/domain/session"

// Repository Persistency interface to serve the Session service
type Repository interface {
	InsertVote(Vote) error
	FindSession(string) (session.Session, error)
}
