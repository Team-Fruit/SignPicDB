package messages

import (
	"net/http"

	"github.com/labstack/echo"
)

type handler struct {
	MessageModel message.MessageModelImpl
}

func NewHandler(m user.MessageModelImpl) *handler {
	return &handler{m}
}

func (h *handler) PostMessage(c echo.Context) error {
	m := new(message.Messge)
	if err = c.Bind(m); err != nil {
		return
	}
	if err = c.Validate(m); err != nil {
		return
	}

	m.IP = c.RealIP()
	m.Message = ""

	if err = h.MessageModel.Push(m); err != nil {
		return
	}

	return c.JSON(http.StatusOK, m)
}
