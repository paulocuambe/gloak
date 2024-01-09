package realm

import (
	"context"
	"database/sql"
	"log"

	"github.com/paulocuambe/gloak/internal/db"
	"github.com/paulocuambe/gloak/internal/models"
)

type service struct {
	store *db.DB
}

func (d *service) DB() *sql.DB {
	return d.store.DB
}

func ProvideService(db *db.DB) *service {
	return &service{db}
}

func (s *service) GetRealmByID(ctx context.Context, id string) (*models.Realm, error) {
	query := "SELECT id, name, created_at, updated_at FROM realms WHERE id = $1"
	if s.store.Cfg.IsSqlite3() {
		query = "SELECT id, name, created_at, updated_at FROM realms WHERE id = ?"
	}

	rs := s.DB().QueryRow(query, id)
	if err := rs.Err(); err != nil {
		return nil, err
	}

	var realm models.Realm
	err := rs.Scan(&realm.Id, &realm.Name, &realm.CreatedAt, &realm.UpdatedAt)

	if err != nil {
		return nil, err
	}

	log.Printf("realm: %#v", realm)
	return &realm, nil
}

func (s *service) Create(ctx context.Context, cmd *models.CreateRealmCommand) (*models.Realm, error) {
	tx, err := s.DB().BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	query := "INSERT INTO realms (id, name) values ($1, $2)"
	if s.store.Cfg.IsSqlite3() {
		query = "INSERT INTO realms (id, name) values (?, ?)"
	}

	log.Println(query)

	rs, err := tx.Exec(query, cmd.Id, cmd.Name)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	log.Println("id: ", id)

	if err != nil {
		return nil, err
	}

	return nil, tx.Commit()
}
