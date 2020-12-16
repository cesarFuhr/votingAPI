package session

import (
	"testing"
	"time"
)

func TestGetExpiration(t *testing.T) {
	s := Session{
		ID:             "andID",
		OriginalAgenda: "originalAgenda",
		Creation:       time.Now(),
		Duration:       time.Minute,
	}
	t.Run("returns the expiration timestamp", func(t *testing.T) {
		got := s.GetExpiration()
		want := s.Creation.Add(s.Duration)

		assertValue(t, got, want)
	})
}
