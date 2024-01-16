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
	log   *log.Logger
}

func (d *service) DB() *sql.DB {
	return d.store.DB
}

func ProvideService(db *db.DB) *service {
	return &service{store: db}
}

func (s *service) GetRealmByID(ctx context.Context, id string) (*models.Realm, error) {
	tx, err := s.DB().Begin()
	if err != nil {
		return nil, err
	}

	realm, err := s.get(tx, id)
	return realm, tx.Commit()
}

func (s *service) Create(ctx context.Context, cmd *models.CreateRealmCommand) (*models.Realm, error) {
	tx, err := s.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = s.insert(tx, cmd)
	if err != nil {
		return nil, err
	}

	r, err := s.get(tx, cmd.Id)
	if err != nil {
		return nil, err
	}

	return r, tx.Commit()
}

func (s *service) insert(tx *sql.Tx, cmd *models.CreateRealmCommand) error {
	query := "INSERT INTO realms (id, name) values ($1, $2)"
	if s.store.Cfg.IsSqlite3() {
		query = "INSERT INTO realms (id, name) values (?, ?)"
	}

	_, err := tx.Exec(query, cmd.Id, cmd.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) get(tx *sql.Tx, id string) (*models.Realm, error) {
	query := "SELECT id, name, created_at, updated_at from realms WHERE id=$1"
	if s.store.Cfg.IsSqlite3() {
		query = "SELECT id, name, created_at, updated_at from realms WHERE id=?"
	}

	log.Println(query)
	r := tx.QueryRow(query, id)
	if err := r.Err(); err != nil {
		return nil, err
	}

	var realm models.Realm

	err := r.Scan(&realm.Id, &realm.Name, &realm.CreatedAt, &realm.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &realm, nil
}

func (s *service) getAll(tx *sql.Tx, id string) ([]*models.Realm, error) {
	query := "SELECT id, name, created_at, updated_at from realms"
	if s.store.Cfg.IsSqlite3() {
		query = "SELECT id, name, created_at, updated_at from realms"
	}

	log.Println(query)
	rs, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	var realms []*models.Realm

	for {
		if !rs.Next() {
			break
		}

		var realm models.Realm
		err := rs.Scan(&realm.Id, &realm.Name, &realm.CreatedAt, &realm.UpdatedAt)
		if err != nil {
			return nil, err
		}

		realms = append(realms, &realm)
	}

	return realms, nil
}
