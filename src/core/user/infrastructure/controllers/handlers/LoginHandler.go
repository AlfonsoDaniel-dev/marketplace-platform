package userHandlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"shopperia/src/common/models"
	"shopperia/src/common/responses"
)

func (H *Handler) Login(c echo.Context) error {
	form := models.Login{}

	err := c.Bind(&form)
	if err != nil {
		response := responses.NewResponse("error", "Bad structured Request", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	token, err := H.Service.LoginUser(form)
	if err != nil {
		response := responses.NewResponse("error", "error while login", err)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := responses.GenerateResponses("ok", "login Success", []any{token})
	return c.JSON(http.StatusOK, response)
}
