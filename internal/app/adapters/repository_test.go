package adapters

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cesarFuhr/votingAPI/internal/app/domain/agenda"
	"github.com/cesarFuhr/votingAPI/internal/app/domain/session"
	"github.com/cesarFuhr/votingAPI/internal/app/domain/vote"
)

type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var agendaMock = agenda.Agenda{
	ID:          "string",
	Description: "string",
}

var sessionMock = session.Session{
	ID:             "string",
	OriginalAgenda: "string",
	Duration:       time.Minute,
	Creation:       time.Now(),
}

var voteMock = vote.Vote{
	AssociateID: "string",
	SessionID:   "string",
	Document:    "12333",
	Vote:        "S",
	Creation:    time.Now(),
}

type loggerStub struct{}

func (l *loggerStub) Info(...interface{}) {}

func TestInsertAgenda(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db, l: &loggerStub{}}
	defer db.Close()

	t.Run("calls db.Exec with the right params", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO agenda").WithArgs(
			agendaMock.ID,
			agendaMock.Description,
			anyTime{},
		)

		repo.InsertAgenda(agendaMock)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectExec("INSERT INTO agendas").WithArgs(
			agendaMock.ID,
			agendaMock.Description,
			anyTime{},
		).WillReturnError(want)

		got := repo.InsertAgenda(agendaMock)

		assertValue(t, got, want)
	})
}

func TestFindAgenda(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db, l: &loggerStub{}}
	defer db.Close()

	t.Run("calls db.QueryRow with the right params", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT id, description
				FROM agendas
				WHERE id`).WithArgs(agendaMock.ID)

		repo.FindAgenda(agendaMock.ID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("returns a complete Agenda object", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "description"}).
			AddRow(agendaMock.ID, agendaMock.Description)
		mock.
			ExpectQuery(`
					SELECT id, description
						FROM agendas
						WHERE id`).
			WithArgs(agendaMock.ID).
			WillReturnRows(rows)

		returned, err := repo.FindAgenda(agendaMock.ID)

		assertValue(t, err, nil)
		if !reflect.DeepEqual(agendaMock, returned) {
			t.Errorf("want %v, got %v", agendaMock, returned)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectQuery(`
				SELECT id, description
					FROM agendas
					WHERE id`).WithArgs(agendaMock.ID).WillReturnError(want)

		_, got := repo.FindAgenda(agendaMock.ID)

		assertValue(t, got, want)
	})

	t.Run("not founding the key, return a Agenda not Found", func(t *testing.T) {
		want := errors.New("Agenda not found")
		mock.ExpectQuery(`
				SELECT id, description
					FROM agendas
					WHERE id`).WithArgs(agendaMock.ID).WillReturnRows(sqlmock.NewRows([]string{}))

		_, got := repo.FindAgenda(agendaMock.ID)

		assertValue(t, got.Error(), want.Error())
	})
}

func TestInsertSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db, l: &loggerStub{}}
	defer db.Close()

	t.Run("calls db.Exec with the right params", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sessions").WithArgs(
			sessionMock.ID,
			sessionMock.OriginalAgenda,
			sessionMock.Duration,
			anyTime{},
		)

		repo.InsertSession(sessionMock)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectExec("INSERT INTO sessions").WithArgs(
			sessionMock.ID,
			sessionMock.OriginalAgenda,
			sessionMock.Duration,
			anyTime{},
		).WillReturnError(want)

		got := repo.InsertSession(sessionMock)

		assertValue(t, got, want)
	})
}

func TestFindSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db, l: &loggerStub{}}
	defer db.Close()

	t.Run("calls db.QueryRow with the right params", func(t *testing.T) {
		mock.ExpectQuery(`
		SELECT id, originalAgenda, duration, creation
				FROM sessions
				WHERE id`).WithArgs(sessionMock.ID)

		repo.FindSession(sessionMock.ID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("returns a complete Session object", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "originalAgenda", "duration", "creation"}).
			AddRow(sessionMock.ID, sessionMock.OriginalAgenda, sessionMock.Duration, sessionMock.Creation)
		mock.
			ExpectQuery(`
					SELECT id, originalAgenda, duration, creation
					FROM sessions
					WHERE id`).
			WithArgs(sessionMock.ID).
			WillReturnRows(rows)

		returned, err := repo.FindSession(sessionMock.ID)

		assertValue(t, err, nil)
		if !reflect.DeepEqual(sessionMock, returned) {
			t.Errorf("want %v, got %v", sessionMock, returned)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectQuery(`
		SELECT id, originalAgenda, duration, creation
		FROM sessions
		WHERE id`).WithArgs(sessionMock.ID).WillReturnError(want)

		_, got := repo.FindSession(sessionMock.ID)

		assertValue(t, got, want)
	})

	t.Run("not founding the key, return a Session not found", func(t *testing.T) {
		want := errors.New("Session not found")
		mock.ExpectQuery(`
				SELECT id, originalAgenda, duration, creation
				FROM sessions
					WHERE id`).WithArgs(sessionMock.ID).WillReturnRows(sqlmock.NewRows([]string{}))

		_, got := repo.FindSession(sessionMock.ID)

		assertValue(t, got.Error(), want.Error())
	})
}

func TestInsertVote(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db, l: &loggerStub{}}
	defer db.Close()

	t.Run("calls db.Exec with the right params", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO votes").WithArgs(
			voteMock.AssociateID,
			voteMock.SessionID,
			voteMock.Document,
			voteMock.Vote,
			anyTime{},
		)

		repo.InsertVote(voteMock)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectExec("INSERT INTO votes").WithArgs(
			voteMock.AssociateID,
			voteMock.SessionID,
			voteMock.Document,
			voteMock.Vote,
			anyTime{},
		).WillReturnError(want)

		got := repo.InsertVote(voteMock)

		assertValue(t, got, want)
	})
}

func TestFindVotes(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db, l: &loggerStub{}}
	defer db.Close()

	t.Run("calls db.Exec with the right params", func(t *testing.T) {
		mock.ExpectQuery("SELECT vote FROM votes WHERE sessionID").
			WithArgs(sessionMock.ID)

		repo.FindVotes(sessionMock)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("returns a list of vote values", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"vote"}).
			AddRow("S").
			AddRow("N").
			AddRow("S").
			AddRow("S")
		mock.
			ExpectQuery(`
					SELECT vote
					FROM votes
					WHERE sessionID`).
			WithArgs(sessionMock.ID).
			WillReturnRows(rows)

		returned, err := repo.FindVotes(sessionMock)

		assertValue(t, err, nil)
		if !reflect.DeepEqual([]string{"S", "N", "S", "S"}, returned) {
			t.Errorf("want %v, got %v", sessionMock, returned)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectQuery("SELECT vote FROM votes WHERE sessionID").
			WithArgs(sessionMock.ID).WillReturnError(want)

		_, got := repo.FindVotes(sessionMock)

		assertValue(t, got, want)
	})
}

func assertValue(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func assertType(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %T, got %T", want, got)
	}
}
