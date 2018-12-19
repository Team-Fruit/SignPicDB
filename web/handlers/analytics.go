package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/Team-Fruit/SignPicDB/web/models"
)

func (h *handler) GetAnalytics(c echo.Context) (err error) {
	var d models.AnalyticsData
	if d, err = h.Model.GetAnalyticsData(); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, d)
}
