package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	server "github.com/cesarFuhr/votingAPI/internal/app"
	"github.com/cesarFuhr/votingAPI/internal/app/adapters"
	"github.com/cesarFuhr/votingAPI/internal/app/domain/agenda"
	"github.com/cesarFuhr/votingAPI/internal/app/domain/session"
	"github.com/cesarFuhr/votingAPI/internal/app/domain/vote"
	"github.com/cesarFuhr/votingAPI/internal/app/ports"
	"github.com/cesarFuhr/votingAPI/internal/pkg/config"
	"github.com/cesarFuhr/votingAPI/internal/pkg/db"
	"github.com/cesarFuhr/votingAPI/internal/pkg/logger"
)

func main() {
	run()
}

func run() {

	cfgSource := getCfgSource()
	cfg, err := config.LoadConfigs(cfgSource)
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)
	db := bootstrapSQLDatabase(cfg)
	httpServer := bootstrapHTTPServer(cfg, db)

	if err := http.ListenAndServe(":"+cfg.Server.Port, httpServer); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

func bootstrapSQLDatabase(cfg config.Config) *sql.DB {
	sqlDB, err := db.NewPGDatabase(db.PGConfigs{
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		User:     cfg.Db.User,
		Password: cfg.Db.Password,
		Dbname:   cfg.Db.Dbname,
		Driver:   cfg.Db.Driver,
	})
	if err != nil {
		panic(err)
	}
	db.RunMigrations(sqlDB)
	return sqlDB
}

func bootstrapHTTPServer(cfg config.Config, sqlDB *sql.DB) server.HTTPServer {
	l := logger.NewLogger()

	sqlRepo := adapters.NewSQLRepository(sqlDB, l)
	mqttPub := adapters.NewMQTTPublisher(cfg.Broker.ConnString, l)

	agendaService := agenda.NewAgendaService(&sqlRepo)
	agendaHandler := ports.NewAgendaHandler(agendaService)

	sessionService := session.NewSessionService(&sqlRepo, &mqttPub)
	sessionHandler := ports.NewSessionHandler(sessionService)
	resultHandler := ports.NewResultHandler(sessionService)

	voteService := vote.NewVoteService(&sqlRepo, &adapters.DocValidator{})
	voteHandler := ports.NewVoteHandler(voteService)

	return server.NewHTTPServer(l, agendaHandler, sessionHandler, voteHandler, resultHandler)
}

func getCfgSource() string {
	var cfgFromEnv bool
	flag.BoolVar(&cfgFromEnv, "e", false, "load config from environment")
	flag.Parse()
	if cfgFromEnv == true {
		return "env"
	}
	return "yaml"
}
