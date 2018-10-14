package main

import (
    "fmt"
    "strings"
    "reflect"
)

const user = `REPLACE INTO user (uuid, username, ip, version_mod, version_mod_mc, version_mod_forge, version_mc, version_forge, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

func (u *User) Push() {
    db.MustExec(user, u.UUID,
                      u.UserName,
                      u.IP, u.VersionMod,
                      u.VersionModMC,
                      u.VersionModForge,
                      u.VersionMC,
                      u.VersionForge,
                      u.Message)
}

func (w *Where) Pull(offset uint64, limit uint64) (u []User, err error) {
    ws := w.toSql()
    if ws != "" {
        err = db.Select(&u, fmt.Sprintf("SELECT * FROM user WHERE %s LIMIT %d,%d", ws, offset, limit))
    } else {
        err = db.Select(&u, fmt.Sprintf("SELECT * FROM user LIMIT %d,%d", offset, limit))
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
            s = append(s, fmt.Sprintf("%s%s'%s'", tag.Get("db"), tag.Get("operator"), fv))
        }
    }

    return strings.Join(s, " AND ")
}

