package adapters

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/agenda"
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

func TestInsertAgenda(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db}
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

func TestSQLFindKey(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLRepository{db: db}
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

	t.Run("not founding the key, return a ErrKeyNotFound", func(t *testing.T) {
		want := errors.New("Agenda not found")
		mock.ExpectQuery(`
				SELECT id, description
					FROM agendas
					WHERE id`).WithArgs(agendaMock.ID).WillReturnRows(sqlmock.NewRows([]string{}))

		_, got := repo.FindAgenda(agendaMock.ID)

		assertValue(t, got.Error(), want.Error())
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
