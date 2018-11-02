package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type (
	User struct {
		UUID            string `db:"uuid" query:"id" validate:"required,len=32"`
		UserName        string `db:"username" query:"name" validate:"required"`
		IP              string `db:"ip"`
		VersionMod      string `db:"version_mod" query:"vmod" validate:"required"`
		VersionModMC    string `db:"version_mod_mc" query:"vmodmc" validate:"required"`
		VersionModForge string `db:"version_mod_forge" query:"vmodforge" validate:"required"`
		VersionMC       string `db:"version_mc" query:"vmc" validate:"required"`
		VersionForge    string `db:"version_forge" query:"vforge" validate:"required"`
		Message         string `db:"message"`
		UpdatedAt       string `db:"updated_at"`
	}
	Where struct {
		UUID              string `query:"id" validate:"omitempty,len=32" db:"uuid" operator:"="`
		UserName          string `query:"name" db:"username" operator:"="`
		IP                string `query:"ip" validate:"omitempty,ip" db:"ip" operator:"="`
		Version_Mod       string `query:"vmod" db:"version_mod" operator:"="`
		Version_Mod_MC    string `query:"vmodmc" db:"version_mod_mc" operator:"="`
		Version_Mod_Forge string `query:"vmodforge" db:"version_mod_forge" operator:"="`
		Version_MC        string `query:"vmc" db:"version_mc" operator:"="`
		Version_Forge     string `query:"vforge" db:"version_forge" operator:"="`
		// Since             string `query:"since" db:"updated_at" operator:">="`
		// Until             string `query:"until" db:"updated_at" operator:"<="`
	}
	Count struct {
		Count uint64
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

var db *sqlx.DB

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	db = sqlx.MustConnect("mysql", "signpic:@tcp(db:3306)/signpic_db")
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", root)
	e.POST("/", root)
	e.GET("/list", list)
	e.GET("/list/count", count)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func root(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}
	if err = c.Validate(u); err != nil {
		return
	}

	u.IP = c.RealIP()
	u.Message = ""

	u.Push()

	return c.JSON(http.StatusOK, u)
}

func list(c echo.Context) (err error) {
	var page, pagesize uint64
	if pagestr := c.QueryParam("page"); pagestr != "" {
		if page, err = strconv.ParseUint(pagestr, 10, 32); err != nil {
			return
		}
		if page < 1 {
			page = 1
		}
	} else {
		page = 1
	}
	if pagesizestr := c.QueryParam("pagesize"); pagesizestr != "" {
		if pagesize, err = strconv.ParseUint(pagesizestr, 10, 32); err != nil {
			return
		}
		if pagesize < 1 {
			pagesize = 1
		}
		if pagesize > 100 {
			pagesize = 100
		}
	} else {
		pagesize = 100
	}

	w := new(Where)
	if err = c.Bind(w); err != nil {
		return
	}
	if err = c.Validate(w); err != nil {
		return
	}

	var l []User
	if l, err = w.Pull(pagesize*(page-1), pagesize); err != nil {
		return
	}
	if len(l) == 0 {
		l = make([]User, 0)
	}

	return c.JSON(http.StatusOK, l)
}

func count(c echo.Context) (err error) {
	w := new(Where)
	if err = c.Bind(w); err != nil {
		return
	}
	if err = c.Validate(w); err != nil {
		return
	}
	var count uint64
	count, err = w.UserCount()
	return c.JSON(http.StatusOK, Count{count})
}
