package main

import (
	"fmt"
	"os"
	"testing"

	server "github.com/cesarFuhr/gocrypto/internal/app"
	"github.com/cesarFuhr/gocrypto/internal/pkg/config"
	"github.com/cesarFuhr/gocrypto/internal/pkg/db"
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

	sqlDB := bootstrapSQLDatabase(cfg)
	defer sqlDB.Close()

	err = db.RunMigrations(sqlDB)
	if err != nil {
		panic(err)
	}

	httpServer = bootstrapHTTPServer(cfg, sqlDB)

	return m.Run()
}
