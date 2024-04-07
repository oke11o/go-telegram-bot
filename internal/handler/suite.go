package handler

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/oke11o/go-telegram-bot/internal/config"
	"github.com/oke11o/go-telegram-bot/internal/repository/sqlite"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	dbx   *sqlx.DB
	repo  *sqlite.Repo
	dbCfg config.SqliteConfig
	cfg   config.Config
}

func (s *Suite) createDB(cfg config.SqliteConfig) (*sqlx.DB, error) {
	db, err := sql.Open("sqlite3", cfg.File)
	if err != nil {
		return nil, fmt.Errorf("sql.Open() err: %w", err)
	}
	dbx := sqlx.NewDb(db, "sqlite3")
	return dbx, nil
}
