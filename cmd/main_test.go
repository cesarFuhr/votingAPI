package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	server "github.com/cesarFuhr/gocrypto/internal/app"
	"github.com/cesarFuhr/gocrypto/internal/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// loads the file driver to migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var httpServer server.HTTPServer

func TestMain(m *testing.M) {
	os.Exit(deferable(m))
}

func deferable(m *testing.M) int {
	cfg, err := config.LoadConfigs("env")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	db := bootstrapSQLDatabase(cfg)
	defer db.Close()

	err = runMigrations(db)
	if err != nil {
		panic(err)
	}

	httpServer = bootstrapHTTPServer(cfg, db)

	return m.Run()
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../internal/pkg/db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
	return nil
}
