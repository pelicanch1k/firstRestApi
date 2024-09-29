package database

import (
	"fmt"
	"log/slog"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pelicanch1k/homework-http/internal/config"
	"github.com/pelicanch1k/homework-http/internal/storage/psql"
)

type ConnectionInfo struct {
	Username string
	DBName   string
	Password string
}

func NewConnectionInfo(config *config.Config) ConnectionInfo {
	cfg := config.Db

	return ConnectionInfo{
		cfg.Username,
		cfg.DBName,
		cfg.Password,
	}
}

func NewPostgresConnection(info ConnectionInfo, log *slog.Logger) (*psql.Storage, error) {
	const op = "pkg.database.NewPostgresConnection"
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", info.Username, info.Password, info.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &psql.Storage{db, log}, nil
}
