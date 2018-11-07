package users

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/Team-Fruit/SignPicDB/models"
)

type (
	count struct {
		count uint
	}

	handler struct {
		UserModel user.UserModelImpl
	}
)

func NewHandler(u user.UserModelImpl) *handler {
	return &handler{u}
}

func (h *handler) GetList(c echo.Context) error {
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

	w := new(user.Where)
	if err = c.Bind(w); err != nil {
		return
	}
	if err = c.Validate(w); err != nil {
		return
	}

	var l []user.User
	if l, err = h.UserModel.Find(w, pagesize*(page-1), pagesize); err != nil {
		return
	}
	if len(l) == 0 {
		l = make([]user.User, 0)
	}

	return c.JSON(http.StatusOK, l)
}

func (h *handler) GetUniqueUserCount(c echo.Context) error {
	w := new(Where)
	if err = c.Bind(w); err != nil {
		return
	}
	if err = c.Validate(w); err != nil {
		return
	}
	if count, err := h.UserModel.CountUniqueUser(w); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, count{count})
}
