package models

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/paulocuambe/gloak/internal/errors"
)

type Realm struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Realm) ToJson() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type RealmService interface {
	GetRealms(context.Context) ([]*Realm, error)
	GetRealmByID(context.Context, string) (*Realm, error)
	Create(context.Context, *CreateRealmCommand) (*Realm, error)
}

type CreateRealmCommand struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (c *CreateRealmCommand) Validate() []*errors.FieldError {
	errs := make([]*errors.FieldError, 0, 0)
	rt := reflect.TypeOf(c)
	if c.Id == "" {
		f, _ := rt.FieldByName("Id")
		n := f.Tag.Get("json")
		errs = append(errs, errors.NewFieldError(c.Id, n, "'id' can't be null or empty"))
	}
	if c.Name == "" {
		errs = append(errs, errors.NewFieldError(c.Name, "", "'name' can't be null or empty"))
	}
	return errs
}
