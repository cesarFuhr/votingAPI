package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// loads the file driver to migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// PGConfigs configuration for a postgres database
type PGConfigs struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	Dbname       string
	MaxOpenConns int
}

// NewPGDatabase Created a connection with the database and returns it
func NewPGDatabase(cfg PGConfigs) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open(cfg.Driver, psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	fmt.Println("Connected to PGDatabase")
	return db, nil
}

// RunMigrations runs the migrations UP
func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../build/internal/pkg/db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			panic(err)
		}
	}
	return nil
}
