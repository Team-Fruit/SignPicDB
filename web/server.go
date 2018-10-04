package main

import (
    "net/http"

    "github.com/labstack/echo"
)

func main() {
    e := echo.New()

    e.Logger.Fatal(e.Start(":8080"))
}

