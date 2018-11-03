package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

const user = `INSERT INTO user (uuid, username, ip, version_mod, version_mod_mc, version_mod_forge, version_mc, version_forge, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), ip = VALUES(ip), version_mod = VALUES(version_mod), version_mod_mc = VALUES(version_mod_mc), version_mod_forge = VALUES(version_mod_forge), version_mc = VALUES(version_mc), version_forge = (version_forge), message = VALUES(message), updated_at = NOW()`

func (u *User) Push() {
	db.MustExec(user, u.UUID,
		u.UserName,
		u.IP,
		u.VersionMod,
		u.VersionModMC,
		u.VersionModForge,
		u.VersionMC,
		u.VersionForge,
		u.Message)
}

func (w *Where) Pull(offset uint64, limit uint64) (u []User, err error) {
	ws := w.toSql()
	if ws != "" {
		var nstmt *sqlx.NamedStmt
		if nstmt, err = db.PrepareNamed(fmt.Sprintf("SELECT * FROM user WHERE %s LIMIT %d,%d", ws, offset, limit)); err != nil {
			return
		}
		nstmt.Select(&u, w)
	} else {
		err = db.Select(&u, fmt.Sprintf("SELECT * FROM user LIMIT %d,%d", offset, limit))
	}
	return
}

func (w *Where) UserCount() (c uint64, err error) {
	ws := w.toSql()
	if ws != "" {
		var nstmt *sqlx.NamedStmt
		if nstmt, err = db.PrepareNamed("SELECT count(uuid) FROM user WHERE "+ ws); err != nil {
			return
		}
		nstmt.Get(&c, w)
	} else {
		err = db.Get(&c, "SELECT count(uuid) FROM user")
	}
	return
}

func (w *Where) toSql() string {
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

