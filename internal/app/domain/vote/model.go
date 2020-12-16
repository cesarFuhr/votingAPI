package vote

import "time"

// Vote Representation of a vote
type Vote struct {
	AssociateID string
	SessionID   string
	Document    string
	Vote        string
	Creation    time.Time
}
