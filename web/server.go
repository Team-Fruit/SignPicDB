package main

import (
    "net/http"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "gopkg.in/go-playground/validator.v9"
)

type (
    User struct {
        UUID            string `query:"id" validate:"required,len=32"`
        UserName        string `query:"name" validate:"required"`
        VersionMod      string `query:"vmod" validate:"required"`
        VersionModMC    string `query:"vmodmc" validate:"required"`
        VersionModForge string `query:"vmodforge" validate:"required"`
        VersionMC       string `query:"vmc" validate:"required"`
        VersionForge    string `query:"vforge" validate:"required"`
    }

    CustomValidator struct {
        validator *validator.Validate
    }
)

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
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
    return c.JSON(http.StatusOK, u)
}

