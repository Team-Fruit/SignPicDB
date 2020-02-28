package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/Team-Fruit/SignPicDB/web/handlers"
	"github.com/Team-Fruit/SignPicDB/web/models"
	"github.com/Team-Fruit/SignPicDB/web/ws"
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
	db := sqlx.MustConnect("mysql", "signpic:password@tcp(db:3306)/signpic_db")
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	validator := validator.New()
	validator.RegisterValidation("mcuuid", uuidValidator)
	e.Validator = &CustomValidator{validator: validator}

	h := handlers.NewHandler(models.NewModel(db))

	e.POST("/msg", h.PutMessage)
	e.GET("/msg", h.PutMessage)
	e.GET("/users", h.GetList)
	e.GET("/users/:id", h.GetUser)
	e.GET("/usercount", h.GetUniqueUserCount)
	e.GET("/playcount", h.GetPlayCount)
	e.GET("/analytics", h.GetAnalytics)
	e.GET("/analytics/transition", h.GetUserTransition)

	hub := ws.NewHub()
	go hub.Run()
	e.GET("/ws", func(c echo.Context) error {
		return ws.ServeWs(hub, c)
	})
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
