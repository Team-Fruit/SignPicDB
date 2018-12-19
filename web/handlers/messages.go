package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/Team-Fruit/SignPicDB/web/models"
	"github.com/Team-Fruit/SignPicDB/web/ws"
)

func (h *handler) PutMessage(c echo.Context) (err error) {
	m := new(models.Message)
	if err = c.Bind(m); err != nil {
		return
	}
	if err = c.Validate(m); err != nil {
		return
	}

	m.IP = c.RealIP()
	m.Message = ""

	if err = h.Model.PutMessage(m); err != nil {
		return
	}

	var a models.AnalyticsData
	if a, err = h.Model.GetAnalyticsData(); err != nil {
		return
	}
	ws.AnalyticsChan <- a

	return c.JSON(http.StatusOK, m)
}
