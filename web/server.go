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
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func uuidValidator(fl validator.FieldLevel) bool {
	str := fl.Field().String();
	if len(str) != 32 {
		return false
	}
	if str == "00000000000000000000000000000000" {
		return false
	}
	return true
}

func main() {
	db := sqlx.MustConnect("mysql", "signpic:@tcp(db:3306)/signpic_db")
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	validator := validator.New()
	validator.RegisterValidation("mcuuid", uuidValidator)
	e.Validator = &CustomValidator{validator: validator}

	uh := users.NewHandler(user.NewUserModel(db))
	mh := messages.NewHandler(user.NewMessageModel(db))
	
	e.POST("/msg", root)
	e.GET("/msg", root)
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

	if err = u.Push(); err != nil {
		return
	}

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
