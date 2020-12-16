package session

import "time"

// Session Representation of a agenda voting session
type Session struct {
	ID             string
	OriginalAgenda string
	Duration       time.Duration
	Creation       time.Time
}

// GetExpiration Returns the datetime the session wil expire
func (s *Session) GetExpiration() time.Time {
	return s.Creation.Add(s.Duration)
}
