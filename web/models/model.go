package models

import (
	"github.com/jmoiron/sqlx"
)

type (
	Model struct {
		db *sqlx.DB
	}

	Database interface {
		PutMessage(m *Message) error
		FindUsers(where *UserWhere, offset uint, limit uint) (user []User, err error)
		CountUniqueUser(where *UserWhere) (count uint, err error)
	}
)

func NewModel(db *sqlx.DB) *Model {
	return &Model{
		db: db,
	}
}
