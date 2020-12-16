package session

import "time"

// Session Representation of a meeting agenda
type Session struct {
	ID             string
	OriginalAgenda string
	Duration       time.Duration
	Creation       time.Time
}
