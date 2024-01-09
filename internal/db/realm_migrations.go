package db

import (
	"context"
	"fmt"
	"log"
)

func (d *DB) runRealmMigrations(ctx context.Context) error {
	if d.Cfg.IsPostgres() {
	} else if d.Cfg.IsSqlite3() {

		tx, err := d.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		rs, err := tx.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL UNIQUE, updated_at DATETIME default CURRENT_TIMESTAMP, created_at DATETIME default CURRENT_TIMESTAMP)", TABLE_NAME_REALMS))
		if err != nil {
			return err
		}

		i, err := rs.RowsAffected()
		log.Printf("%#v, %#v", i, err)
		i2, err := rs.LastInsertId()
		log.Printf("%#v, %#v", i2, err)

		if err = tx.Commit(); err != nil {
			log.Println(err)
			return err
		}

	}
	return nil
}
