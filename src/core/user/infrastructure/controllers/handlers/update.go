package userHandlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"shopperia/src/auth"
	"shopperia/src/common/responses"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

func (H *Handler) UpdateUserName(c echo.Context) error {
	form := UserDTO.UpdateUserName{}

	err := c.Bind(&form)

	if err != nil {
		fmt.Println(err)
		response := responses.NewResponse("error", "request bad structured", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	token := c.Request().Header.Get("authorization")

	email, err := auth.GetFromToken(token, "email")
	if err != nil {
		response := responses.NewResponse("error", "invalid token", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	form.Common.Email = email

	err = H.Service.UpdateUserName(form)
	if err != nil {
		response := responses.NewResponse("error", "update user name failed", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := responses.GenerateResponses("ok", "username updated successfully", nil)
	return c.JSON(http.StatusOK, response)
}
