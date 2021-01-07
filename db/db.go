package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "gopkg.in/inconshreveable/log15.v2"
	"os"
)

var DB *pgxpool.Pool

func init() {
	logger := log15adapter.NewLogger(log.New("module", "pgx"))
	maxConns := "10"
	if os.Getenv("DB_MAX_CONNS") != "" {
		maxConns = os.Getenv("DB_MAX_CONNS")
	}
	databaseUrl := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		maxConns)
	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Crit("Unable to parse database settings", "error", err)
		os.Exit(1)
	}

	poolConfig.ConnConfig.Logger = logger

	DB, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
}
