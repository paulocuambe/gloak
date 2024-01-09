package db

import (
	"context"
	"log"
)

const (
	TABLE_NAME_REALMS = "realms"
)

func (d *DB) RunMigrations(ctx context.Context) error {
	log.Println("run migrations")
	err := d.runRealmMigrations(ctx)
	if err != nil {
		log.Println("ooohhh")
		log.Println(err)
		return err
	}
	log.Println("finished running migrations")

	return nil
}
