package adapters

import (
	"database/sql"
	"errors"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/agenda"

	// Loading the pq driver
	_ "github.com/lib/pq"
)

// NewSQLRepository returns a new sql repository instance
func NewSQLRepository(db *sql.DB) SQLRepository {
	return SQLRepository{db: db}
}

// SQLRepository sql database persistency
type SQLRepository struct {
	db *sql.DB
}

var findAgendaStatement = `
	SELECT id, description
		FROM agendas
		WHERE id = $1`

// FindAgenda finds and returns the requested key
func (r *SQLRepository) FindAgenda(id string) (agenda.Agenda, error) {
	row := r.db.QueryRow(findAgendaStatement, id)

	var a agenda.Agenda

	switch err := row.Scan(&a.ID, &a.Description); err {
	case nil:
		return a, nil
	case sql.ErrNoRows:
		return agenda.Agenda{}, errors.New("Agenda not found")
	default:
		return agenda.Agenda{}, err
	}
}

var insertAgendaStatement = `
	INSERT INTO agendas (id, description, creation)
		VALUES ($1, $2, $3)`

// InsertAgenda Inserts an agenda into the repository
func (r *SQLRepository) InsertAgenda(a agenda.Agenda) error {
	_, err := r.db.Exec(
		insertAgendaStatement,
		a.ID,
		a.Description,
		time.Now(),
	)
	return err
}
