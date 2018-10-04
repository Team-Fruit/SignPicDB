package main

import (
    "net/http"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

type User struct {
    UUID string `query:"id"`
    UserName string `query:"name"`
    VersionMod string `query:"vmod"`
    VersionModMC string `query:"vmodmc"`
    VersionModForge string `query:"vmodforge"`
    VersionMC string `query:"vmc"`
    VersionForge string `query:"vforge"`
}

func main() {
    e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    e.GET("/", root)

    e.Logger.Fatal(e.Start(":8080"))
}

func root(c echo.Context) error {
    u := new(User)
    if err := c.Bind(u); err != nil {
        return err
    }
    return c.JSON(http.StatusOK, u)
}

