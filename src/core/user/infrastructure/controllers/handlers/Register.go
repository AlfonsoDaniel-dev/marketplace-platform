package userHandlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"shopperia/src/common/models"
	"shopperia/src/common/responses"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

func (H *Handler) Register(c echo.Context) error {
	RegisterForm := UserDTO.RegisterDTO{}

	err := c.Bind(&RegisterForm)
	if err != nil {
		response := responses.NewResponse("error", "request bad structured", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if RegisterForm.Age < 18 {
		response := responses.NewResponse("error", "Age must be greater than 18", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	err = H.Service.Register(RegisterForm)
	if err != nil {
		fmt.Println(err)
		response := responses.NewResponse("error", "error while register", err)
		return c.JSON(http.StatusInternalServerError, response)
	}

	loginModel := models.Login{
		Email:    RegisterForm.Email,
		Password: RegisterForm.Password,
	}

	token, err := H.Service.Login(loginModel)
	if err != nil {
		fmt.Println(err)
		response := responses.NewResponse("error", "register was success but it was an error while generating Token for auto login", err)
		return c.JSON(http.StatusInternalServerError, response)
	}

	data := []any{token, RegisterForm}

	response := responses.GenerateResponses("ok", "Register and Login success", data)
	return c.JSON(http.StatusOK, response)
}
