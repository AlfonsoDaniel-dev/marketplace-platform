package userHandlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"shopperia/src/auth"
	"shopperia/src/common/responses"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

func (H *Handler) UploadAddress(c echo.Context) error {
	form := UserDTO.UploadAddressForm{}

	err := c.Bind(&form)
	if err != nil {
		response := responses.NewResponse("error", "request bad structured", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	token := c.Request().Header.Get("Authorization")

	email, err := auth.GetEmailFromToken(token)
	if err != nil {
		fmt.Println(err)
		response := responses.NewResponse("error", "cannot get email", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	err = H.Service.UploadAddress(email, form)
	if err != nil {
		response := responses.NewResponse("error", "cannot upload address", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	responses := responses.GenerateResponses("ok", "Address Upload", []any{form})
	return c.JSON(http.StatusOK, responses)
}
