package db

import (
	"context"
	"fmt"
	"log"
)

const (
	TABLE_NAME_REALMS = "realms"
)

func (d *DB) RunMigrations(ctx context.Context) error {
	log.Println("running migrations")
	err := d.runRealmMigrations(ctx)
	if err != nil {
		return fmt.Errorf("realm migrations: %w", err)
	}
	log.Println("finished running migrations")

	return nil
}
