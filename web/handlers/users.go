package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/Team-Fruit/SignPicDB/web/models"
)

type count struct {
	Count uint `json:"count"`
}

func (h *handler) GetList(c echo.Context) (err error) {
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

	w := models.UserWhere{}
	if err = c.Bind(&w); err != nil {
		return
	}
	if err = c.Validate(&w); err != nil {
		return
	}

	var l []models.User
	if l, err = h.Model.FindUsers(w, uint(pagesize*(page-1)), uint(pagesize)); err != nil {
		return
	}
	if len(l) == 0 {
		l = make([]models.User, 0)
	}

	return c.JSON(http.StatusOK, l)
}

func (h *handler) GetUniqueUserCount(c echo.Context) (err error) {
	w := models.UserWhere{}
	if err = c.Bind(&w); err != nil {
		return
	}
	if err = c.Validate(&w); err != nil {
		return
	}
	
	var cnt uint
	if cnt, err = h.Model.CountUniqueUser(w); err != nil {
		return
	}
	return c.JSON(http.StatusOK, count{cnt})
}
