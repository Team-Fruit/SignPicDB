package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type (
	User struct {
		UUID         string `json:"uuid" db:"uuid"`
		UserName     string `json:"username" db:"username"`
		IP           string `json:"-" db:"ip"`
		Message      string `json:"message" db:"message"`
		CreatedAt    string `json:"created_at" db:"created_at"`
		UpdatedAt    string `json:"updated_at" db:"updated_at"`
		UpdatedCount uint   `json:"updated_count" db:"updated_count"`
	}

	UserWhere struct {
		UUID     string `db:"uuid" query:"id" validate:"omitempty,mcuuid" operator:"="`
		UserName string `db:"username" query:"name" operator:"="`
		IP       string `db:"ip" query:"ip" validate:"omitempty,ip" operator:"="`
	}

	UserVersion struct {
		UUID            string `json:"-" db:"uuid"`
		VersionMod      string `json:"version_mod" db:"version_mod"`
		VersiomModMC    string `json:"version_mod_mc" db:"version_mod_mc"`
		VersionModForge string `json:"version_mod_forge" db:"version_mod_forge"`
		VersionMC       string `json:"version_mc" db:"version_mc"`
		VersionForge    string `json:"version_forge" db:"version_forge"`
		CreatedAt       string `json:"created_at" db:"created_at"`
		UpdatedAt       string `json:"updated_at" db:"updated_at"`
		UpdatedCount    string `json:"updated_count" db:"updated_count"`
	}

	UserData struct {
		User        *User          `json:"user"`
		UserVersion *[]UserVersion `json:"user_version"`
	}
)


func (u *Model) FindUsers(where UserWhere, offset uint, limit uint) (user []User, err error) {
	ws := where.toSql()
	if ws != "" {
		var nstmt *sqlx.NamedStmt
		if nstmt, err = u.db.PrepareNamed(fmt.Sprintf("SELECT * FROM user WHERE %s LIMIT %d,%d", ws, offset, limit)); err != nil {
			return
		}
		nstmt.Select(&user, where)
	} else {
		err = u.db.Select(&user, fmt.Sprintf("SELECT * FROM user LIMIT %d,%d", offset, limit))
	}
	return
}

func (u *Model) CountUniqueUser(where UserWhere) (count uint, err error) {
	ws := where.toSql()
	if ws != "" {
		var nstmt *sqlx.NamedStmt
		if nstmt, err = u.db.PrepareNamed("SELECT count(uuid) FROM user WHERE "+ ws); err != nil {
			return
		}
		nstmt.Get(&count, where)
	} else {
		err = u.db.Get(&count, "SELECT count(uuid) FROM user")
	}
	return
}

func (u *Model) SumPlayCount() (count uint, err error) {
	err = u.db.Get(&count, "SELECT SUM(updated_count) FROM user")
	return
}

func (m *Model) GetUserData(id string) (userdata UserData, err error) {
	u := User{}
	if err = m.db.Get(&u, "SELECT * FROM user WHERE uuid=? OR username=? ORDER BY updated_at DESC LIMIT 1", id, id); err != nil {
		return
	}
	uv := []UserVersion{}
	if err = m.db.Select(&uv, "SELECT * FROM user__version_mc__version_mod WHERE uuid=?", u.UUID); err != nil {
		return
	}
	userdata = UserData{User: &u, UserVersion: &uv}
	return
}

func (w *UserWhere) toSql() string {
	s := []string{}

	val := reflect.ValueOf(w).Elem()
	for i := 0; i < val.NumField(); i++ {
		vf := val.Field(i)
		tf := val.Type().Field(i)
		tag := tf.Tag
		fv := vf.Interface().(string)
		if fv != "" {
			s = append(s, fmt.Sprintf("%s%s:%s", tag.Get("db"), tag.Get("operator"), strings.ToLower(tf.Name)))
		}
	}

	return strings.Join(s, " AND ")
}
