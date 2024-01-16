package db

import (
	"context"
	"fmt"
)

func (d *DB) runRealmMigrations(ctx context.Context) error {
	if d.Cfg.IsPostgres() {
		tx, err := d.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		query := `CREATE TABLE IF NOT EXISTS %v 
		(id VARCHAR(255) PRIMARY KEY, 
		name VARCHAR(255) NOT NULL UNIQUE, 
		updated_at timestamp default CURRENT_TIMESTAMP,
		created_at timestamp default CURRENT_TIMESTAMP)`

		_, err = tx.Exec(fmt.Sprintf(query, TABLE_NAME_REALMS))
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	} else if d.Cfg.IsSqlite3() {
		tx, err := d.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		query := `CREATE TABLE IF NOT EXISTS %v 
		(id VARCHAR(255) PRIMARY KEY, 
		name VARCHAR(255) NOT NULL UNIQUE, 
		updated_at DATETIME default CURRENT_TIMESTAMP,
		created_at DATETIME default CURRENT_TIMESTAMP)`

		_, err = tx.Exec(fmt.Sprintf(query, TABLE_NAME_REALMS))
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

	}
	return nil
}
