package models

import (
	"context"
	"time"
)

type Realm struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateRealmCommand struct {
	Id   string
	Name string
}

type RealmService interface {
	GetRealmByID(context.Context, string) (*Realm, error)
	Create(context.Context, *CreateRealmCommand) (*Realm, error)
}
