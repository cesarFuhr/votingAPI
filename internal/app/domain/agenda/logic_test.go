package agenda

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type AgendaRepoStub struct {
	store map[string]Agenda
}

func (r *AgendaRepoStub) FindAgenda(ID string) (Agenda, error) {
	key, ok := r.store[ID]
	if ok == false {
		return Agenda{}, errors.New("Agenda not found")
	}
	return key, nil
}

func (r *AgendaRepoStub) InsertAgenda(a Agenda) error {
	if a.Description == "error" {
		return errors.New("ops, there was an error")
	}
	r.store[a.ID] = a
	return nil
}

func TestCreateAgenda(t *testing.T) {
	store := map[string]Agenda{}
	repo := AgendaRepoStub{store}
	service := agendaService{&repo}
	t.Run("Returns an agenda", func(t *testing.T) {
		description := "uma descricao da pauta"
		got, _ := service.CreateAgenda(description)
		want := Agenda{}

		assertType(t, got, want)
		assertType(t, got.ID, want.ID)
		assertString(t, got.Description, description)
	})
	t.Run("Returns the error if there was any error", func(t *testing.T) {
		_, got := service.CreateAgenda("error")
		want := errors.New("error")

		assertType(t, got, want)
	})
}

func TestFindAgenda(t *testing.T) {
	store := map[string]Agenda{}
	repo := AgendaRepoStub{store}
	service := agendaService{&repo}
	t.Run("Returns an agenda", func(t *testing.T) {
		want, _ := service.CreateAgenda("description")

		got, _ := service.FindAgenda(want.ID)

		assertType(t, got, want)
		assertString(t, got.ID, want.ID)
		assertType(t, got.Description, want.Description)
	})
	t.Run("Returns an error if there was an error", func(t *testing.T) {
		_, err := service.FindAgenda("notFound")
		want := errors.New("Agenda not found")

		assertType(t, err, want)
		assertString(t, err.Error(), want.Error())
	})
}

func assertType(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("got %T want %T", got, want)
	}
}

func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertTime(t *testing.T, got, want time.Time) {
	t.Helper()
	if got.Round(time.Second) != want.Round(time.Second) {
		t.Errorf("got %v want %v", got, want)
	}
}
