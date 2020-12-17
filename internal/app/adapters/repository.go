package adapters

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/agenda"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/session"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/vote"
	"github.com/cesarFuhr/gocrypto/internal/pkg/logger"

	// Loading the pq driver
	_ "github.com/lib/pq"
)

// NewSQLRepository returns a new sql repository instance
func NewSQLRepository(db *sql.DB, l logger.Logger) SQLRepository {
	return SQLRepository{db: db, l: l}
}

// SQLRepository sql database persistency
type SQLRepository struct {
	db *sql.DB
	l  logger.Logger
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
		r.l.Info(err.Error(), id)
		return agenda.Agenda{}, errors.New("Agenda not found")
	default:
		r.l.Info(err.Error(), id)
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

var findSessionStatement = `
SELECT id, originalAgenda, duration, creation
FROM sessions
WHERE id = $1`

// FindSession finds and returns the requested key
func (r *SQLRepository) FindSession(id string) (session.Session, error) {
	row := r.db.QueryRow(findSessionStatement, id)

	var s session.Session

	switch err := row.Scan(&s.ID, &s.OriginalAgenda, &s.Duration, &s.Creation); err {
	case nil:
		return s, nil
	case sql.ErrNoRows:
		r.l.Info(err.Error(), id)
		return session.Session{}, errors.New("Session not found")
	default:
		r.l.Info(err.Error(), id)
		return session.Session{}, err
	}
}

var insertSessionStatement = `
	INSERT INTO sessions (id, originalAgenda, duration, creation)
		VALUES ($1, $2, $3, $4)`

// InsertSession Inserts an session into the repository
func (r *SQLRepository) InsertSession(s session.Session) error {
	_, err := r.db.Exec(
		insertSessionStatement,
		s.ID,
		s.OriginalAgenda,
		s.Duration,
		s.Creation,
	)
	return err
}

var insertVoteStatement = `
	INSERT INTO votes (associateID, sessionID, document, vote, creation)
		VALUES ($1, $2, $3, $4, $5)`

// InsertVote Inserts a vote into the repository
func (r *SQLRepository) InsertVote(v vote.Vote) error {
	_, err := r.db.Exec(
		insertVoteStatement,
		v.AssociateID,
		v.SessionID,
		v.Document,
		v.Vote,
		v.Creation,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			r.l.Info(err.Error(), v.AssociateID)
			return vote.ErrDuplicateVote
		}
		return err
	}
	return nil
}

var findVotesStatement = `
	SELECT vote
		FROM votes
		WHERE sessionID = $1`

// FindVotes Finds all votes by sessionID
func (r *SQLRepository) FindVotes(s session.Session) ([]string, error) {
	rows, err := r.db.Query(findVotesStatement, s.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var voteValues []string
	for rows.Next() {
		var v string
		err := rows.Scan(&v)
		if err != nil {
			r.l.Info(err.Error(), s.ID)
			return nil, err
		}
		voteValues = append(voteValues, v)
	}
	return voteValues, nil
}
