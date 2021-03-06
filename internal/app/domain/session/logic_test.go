package session

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

type ClockStub struct {
	RightNow time.Time
}

func (c ClockStub) Now() time.Time {
	return c.RightNow
}

type SessionRepoStub struct {
	store map[string]Session
}

func (r *SessionRepoStub) FindSession(ID string) (Session, error) {
	session, ok := r.store[ID]
	if ok == false {
		return Session{}, errors.New("Session not found")
	}
	return session, nil
}

func (r *SessionRepoStub) InsertSession(s Session) error {
	if s.OriginalAgenda == "error" {
		return errors.New("ops, there was an error")
	}
	r.store[s.ID] = s
	return nil
}

func (r *SessionRepoStub) FindVotes(s Session) ([]string, error) {
	if s.OriginalAgenda == "error" {
		return []string{}, errors.New("ops, there was an error")
	}
	return []string{"S", "N", "S", "N", "S"}, nil
}

type PublisherStub struct {
	CalledWith []interface{}
	wg         sync.WaitGroup
}

func (p *PublisherStub) PublishResult(r Result) error {
	return nil
}

func TestCreateSession(t *testing.T) {
	now := time.Now()
	store := map[string]Session{}
	clockStub := ClockStub{RightNow: now}
	pubStub := PublisherStub{CalledWith: []interface{}{}}
	repo := SessionRepoStub{store}
	service := sessionService{&repo, &pubStub, &clockStub}
	t.Run("Returns an session", func(t *testing.T) {
		agendaID := "anID"
		duration := time.Duration(time.Minute) * 5
		got, _ := service.CreateSession(agendaID, duration)
		want := Session{}

		assertType(t, got, want)
		assertType(t, got.Creation, time.Now())
		assertType(t, got.ID, "anotherID")
		assertValue(t, got.OriginalAgenda, agendaID)
		assertValue(t, got.Duration, duration)
		assertValue(t, got.Creation, now)
	})
	t.Run("If informed duration is zero should assume 1 minute", func(t *testing.T) {
		agendaID := "anID"
		duration := time.Duration(time.Minute) * 0
		got, _ := service.CreateSession(agendaID, duration)
		want := time.Minute

		assertValue(t, got.Duration, want)
	})
	t.Run("Returns the error if there was any error", func(t *testing.T) {
		_, got := service.CreateSession("error", time.Duration(time.Minute))
		want := errors.New("ops, there was an error")

		assertType(t, got, want)
	})
}

func TestFindSession(t *testing.T) {
	now := time.Now()
	store := map[string]Session{}
	clockStub := ClockStub{RightNow: now}
	pubStub := PublisherStub{CalledWith: []interface{}{}}
	repo := SessionRepoStub{store}
	service := sessionService{&repo, &pubStub, &clockStub}
	t.Run("Returns an session", func(t *testing.T) {
		agendaID := "anID"
		duration := time.Duration(time.Minute) & 5
		s, _ := service.CreateSession(agendaID, duration)

		got, _ := service.FindSession(s.ID)
		want := Session{}

		assertType(t, got, want)
		assertType(t, got.Creation, time.Now())
		assertType(t, got.ID, "anotherID")
		assertType(t, got.Duration, time.Duration(time.Second))
		assertType(t, got.Creation, now)
		assertValue(t, got.OriginalAgenda, s.OriginalAgenda)
	})
	t.Run("Returns an error if there was an error", func(t *testing.T) {
		_, err := service.FindSession("notFound")
		want := errors.New("Session not found")

		assertType(t, err, want)
		assertValue(t, err.Error(), want.Error())
	})
}

func TestResult(t *testing.T) {
	now := time.Now()
	store := map[string]Session{}
	clockStub := ClockStub{RightNow: now}
	pubStub := PublisherStub{CalledWith: []interface{}{}}
	repo := SessionRepoStub{store}
	service := sessionService{&repo, &pubStub, &clockStub}
	t.Run("Returns an result", func(t *testing.T) {
		agendaID := "anID"
		duration := time.Duration(time.Minute) & 5
		s, _ := service.CreateSession(agendaID, duration)

		got, _ := service.Result(s.ID)
		want := Result{}

		assertType(t, got, want)
		assertValue(t, got.ID, s.ID)
		assertValue(t, got.OriginalAgenda, s.OriginalAgenda)
		assertType(t, got.Count, want.Count)
		assertType(t, got.Closed, want.Closed)
	})
	t.Run("Returns an result closed result if is session is expired", func(t *testing.T) {
		agendaID := "anID"
		duration := time.Duration(time.Minute) & 5
		s, _ := service.CreateSession(agendaID, duration)

		clockStub.RightNow = now.Add(time.Duration(time.Hour))
		got, _ := service.Result(s.ID)
		want := true

		assertValue(t, got.Closed, want)
	})
	t.Run("Returns an result not closed result if is session is not expired", func(t *testing.T) {
		agendaID := "anID"
		duration := time.Duration(time.Minute) & 5
		s, _ := service.CreateSession(agendaID, duration)

		clockStub.RightNow = now.Add(-time.Duration(time.Hour))
		got, _ := service.Result(s.ID)
		want := false

		assertValue(t, got.Closed, want)
	})
	t.Run("Returns count of the votes", func(t *testing.T) {
		agendaID := "anID"
		duration := time.Duration(time.Minute) & 5
		s, _ := service.CreateSession(agendaID, duration)

		clockStub.RightNow = now.Add(-time.Duration(time.Hour))
		got, _ := service.Result(s.ID)
		want := Count{3, 2}

		assertValue(t, got.Count, want)
	})
	t.Run("Returns an error if there was an error", func(t *testing.T) {
		_, err := service.FindSession("notFound")
		want := errors.New("Session not found")

		assertType(t, err, want)
		assertValue(t, err.Error(), want.Error())
	})
}

func assertType(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("got %T want %T", got, want)
	}
}

func assertValue(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertTime(t *testing.T, got, want time.Time) {
	t.Helper()
	if got.Round(time.Second) != want.Round(time.Second) {
		t.Errorf("got %v want %v", got, want)
	}
}
