package userHandlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	UserServices "shopperia/src/core/user/app/Services"
)

type Handler struct {
	Service UserServices.ServiceInterface
}

func NewHandler(service UserServices.ServiceInterface) *Handler {
	return &Handler{
		Service: service,
	}
}

func (H *Handler) HelloWorld(c echo.Context) error {

	text := H.Service.Hello()
	return c.String(http.StatusOK, text)
}
