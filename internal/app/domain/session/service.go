package session

import "time"

// Service describes the agenda service interface
type Service interface {
	CreateSession(string, time.Duration) (Session, error)
	FindSession(string) (Session, error)
}
