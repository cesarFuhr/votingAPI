package vote

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/session"
)

type VoteRepoStub struct {
	sessionStore map[string]session.Session
	voteStore    map[string]Vote
}

func (r *VoteRepoStub) FindSession(ID string) (session.Session, error) {
	s, ok := r.sessionStore[ID]
	if ok == false {
		return session.Session{}, errors.New("Session not found")
	}
	return s, nil
}

func (r *VoteRepoStub) InsertVote(v Vote) error {
	if v.AssociateID == "error" {
		return errors.New("ops, there was an error")
	}
	if _, ok := r.voteStore[v.AssociateID]; ok == true {
		return ErrDuplicateVote
	}
	r.voteStore[v.AssociateID] = v
	return nil
}

func TestCreateVote(t *testing.T) {
	sStore := map[string]session.Session{
		"sessionID": {},
	}
	vStore := map[string]Vote{
		"existing": {},
	}
	repo := VoteRepoStub{sStore, vStore}
	service := voteService{&repo}
	t.Run("Returns an vote", func(t *testing.T) {
		associateID := "anID"
		sessionID := "sessionID"
		document := "01791229005"
		got, _ := service.CreateVote(associateID, sessionID, document)
		want := Vote{}

		assertType(t, got, want)
		assertType(t, got.AssociateID, associateID)
		assertType(t, got.Creation, time.Now())
		assertValue(t, got.SessionID, sessionID)
		assertValue(t, got.Document, document)
	})
	t.Run("Returns an Session not found error if it does not exists", func(t *testing.T) {
		associateID := "anID"
		sessionID := "notFound"
		document := "01791229005"
		_, got := service.CreateVote(associateID, sessionID, document)
		want := errors.New("Session not found")

		assertValue(t, got.Error(), want.Error())
	})
	t.Run("Returns an Duplicate Vote error if its duplicate", func(t *testing.T) {
		associateID := "existing"
		sessionID := "sessionID"
		document := "01791229005"
		_, got := service.CreateVote(associateID, sessionID, document)
		want := ErrDuplicateVote

		assertValue(t, got.Error(), want.Error())
	})
	t.Run("Returns the error if there was any error", func(t *testing.T) {
		associateID := "error"
		sessionID := "sessionID"
		document := "01791229005"
		_, got := service.CreateVote(associateID, sessionID, document)
		want := errors.New("ops, there was an error")

		assertValue(t, got.Error(), want.Error())
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