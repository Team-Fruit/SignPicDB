package handlers

import (
	"github.com/Team-Fruit/SignPicDB/web/models"
)

type handler struct {
	Model models.Database
}

func NewHandler(d models.Database) *handler {
	return &handler{
		Model: d,
	}
}
