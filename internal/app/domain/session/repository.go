package session

// Repository Persistency interface to serve the Session service
type Repository interface {
	FindSession(string) (Session, error)
	InsertSession(Session) error
	FindVotes(Session) ([]string, error)
}
