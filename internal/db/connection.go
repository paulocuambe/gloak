package db

import (
	"database/sql"

	"github.com/paulocuambe/gloak/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Cfg *config.DatabaseConfig
	DB  *sql.DB
}

func initiateConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver.GetName(), cfg.DSN())

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ProvideDBConnection(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := initiateConnection(cfg)

	if err != nil {
		return nil, err
	}

	return &DB{DB: db, Cfg: cfg}, err
}
