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
		FindUsers(where UserWhere, offset uint, limit uint) (user []User, err error)
		CountUniqueUser(where UserWhere) (count uint, err error)
		SumPlayCount() (count uint, err error)
		GetUserData(id string) (userdata UserData, err error)
		GetPlayCount() (c uint, err error)
		GetUserCount() (c uint, err error)
		GetMostPlayedMCVersion() (v string, err error)
		GetMostPlayedModVersion() (v string, err error)
		GetAnalyticsData() (d AnalyticsData, err error)
	}
)

func NewModel(db *sqlx.DB) *Model {
	return &Model{
		db: db,
	}
}
