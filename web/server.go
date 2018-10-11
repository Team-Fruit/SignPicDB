package main

import (
    "log"
    "time"
    "net/http"
    "database/sql"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "gopkg.in/go-playground/validator.v9"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

type (
    User struct {
        UUID            string         `db:"uuid" query:"id" validate:"required,len=32"`
        UserName        string         `db:"username" query:"name" validate:"required"`
        IP              string         `db:"ip"`
        VersionMod      string         `db:"version_mod" query:"vmod" validate:"required"`
        VersionModMC    string         `db:"version_mod_mc" query:"vmodmc" validate:"required"`
        VersionModForge string         `db:"version_mod_forge" query:"vmodforge" validate:"required"`
        VersionMC       string         `db:"version_mc" query:"vmc" validate:"required"`
        VersionForge    string         `db:"version_forge" query:"vforge" validate:"required"`
        Message         sql.NullString `db:"message"`
        UpdatedAt       string         `db:"updated_at"`
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
    var err error
    db, err = sqlx.Connect("mysql", "signpic:@tcp(db:3306)/signpic_db")
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()

    e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    e.Validator = &CustomValidator{validator: validator.New()}

    e.GET("/", root)

    e.Logger.Fatal(e.Start(":8080"))
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
    u.UpdatedAt = time.Now().Format(time.RFC3339)
    
    u.AddUser()

    return c.JSON(http.StatusOK, u)
}

