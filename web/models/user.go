package user

import (
	"github.com/jmoiron/sqlx"
)

type (
	UserModelImpl interface {
		Find(where UserWhere, offset uint, limit uint) (user []User, err error)
		CountUniqueUser(where UserWhere) (count uint, err error)
	}

	User struct {
		UUID         string `json:"uuid" db:"uuid"`
		UserName     string `json:"username" db:"username"`
		IP           string `json:"ip" db:"ip"`
		Message      string `json:"message" db:"message"`
		CreatedAt    string `json:"created_at" db:"created_at"`
		UpdatedAt    string `json:"updated_at" db:"updated_at"`
		UpdatedCount uint   `json:"updated_count" db:"updated_count"`
	}

	UserWhere struct {
		UUID     string `db:"uuid" operator:"="`
		UserName string `db:"username" operator:"="`
		IP       string `db:"ip" operator:"="`
	}

	UserModel struct {
		db *sqlx.DB
	}
)

func NewUserModel(db *sqlx.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (u *UserModel) Find(where UserWhere, offset uint, limit uint) (user []User, err error) {
	ws := where.toSql()
	if ws != "" {
		var nstmt *sqlx.NamedStmt
		if nstmt, err = u.db.PrepareNamed(fmt.Sprintf("SELECT * FROM user WHERE %s LIMIT %d,%d", ws, offset, limit)); err != nil {
			return
		}
		nstmt.Select(&u, w)
	} else {
		err = u.db.Select(&u, fmt.Sprintf("SELECT * FROM user LIMIT %d,%d", offset, limit))
	}
	return
}

func (u *UserModel) CountUniqueUser(where UserWhere) (count uint, err error) {
	ws := where.toSql()
	if ws != "" {
		var nstmt *sqlx.NamedStmt
		if nstmt, err = u.db.PrepareNamed("SELECT count(uuid) FROM user WHERE "+ ws); err != nil {
			return
		}
		nstmt.Get(&c, w)
	} else {
		err = u.db.Get(&c, "SELECT count(uuid) FROM user")
	}
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
