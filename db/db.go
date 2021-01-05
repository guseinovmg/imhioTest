package db

import (
	"context"
	"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "gopkg.in/inconshreveable/log15.v2"
	"os"
)

var DB *pgxpool.Pool

func init() {
	logger := log15adapter.NewLogger(log.New("module", "pgx"))

	poolConfig, err := pgxpool.ParseConfig(`postgresql://postgres:password@localhost:5432/articles`)
	if err != nil {
		log.Crit("Unable to parse DATABASE_URL", "error", err)
		os.Exit(1)
	}

	poolConfig.ConnConfig.Logger = logger

	DB, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
}
