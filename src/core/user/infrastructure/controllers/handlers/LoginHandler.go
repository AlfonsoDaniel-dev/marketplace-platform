package userHandlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"shopperia/src/auth"
	"shopperia/src/common/models"
	"shopperia/src/common/responses"
)

var loginStatus string

func (H *Handler) UserLogin(c echo.Context) error {
	form := models.Login{}

	err := c.Bind(&form)
	if err != nil {
		response := responses.NewResponse("error", "Bad structured Request", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	HasTSV, err := H.Service.CheckTwoStepsVerification(form.Email)
	if err != nil {
		fmt.Println(err)
		response := responses.NewResponse("error", "error while checking TSV", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	if !HasTSV {
		token, err := H.Service.LoginUser(form)
		if err != nil {
			response := responses.NewResponse("error", "error while login", err)
			return c.JSON(http.StatusBadRequest, response)
		}
		response := responses.GenerateResponses("ok", "login Success", []any{token})
		return c.JSON(http.StatusOK, response)
	}

	_, err = H.Service.SendLoginConfirmation(form.Email)
	if err != nil {
		response := responses.NewResponse("error", "error while sending login confirmation", err)
		return c.JSON(http.StatusInternalServerError, response)
	}
	okchan := make(chan string)

	e := c.Echo()

	go e.GET("/api/user/login/confirm/:email/:token", func(c echo.Context) error {
		email := c.Param("email")
		accessToken := c.Param("token")

		ok, err := H.Service.CheckToken(email, accessToken)
		if err != nil {
			fmt.Println(err)
			okchan <- ""
			response := responses.NewResponse("error", "error while checking token", err)
			return c.JSON(http.StatusInternalServerError, response)
		}

		if !ok {
			okchan <- ""
			response := responses.NewResponse("error", "error while checking token", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}

		model := models.Login{
			Email:    email,
			Password: "",
		}

		JwtToken, err := auth.GenerateToken(model, "", false)
		if err != nil {
			okchan <- ""
			response := responses.NewResponse("error", "error while login", err)
			return c.JSON(http.StatusInternalServerError, response)
		}

		okchan <- JwtToken

		err = H.Service.CleanToken(email)
		if err != nil {
			response := responses.NewResponse("error", "error while login", err)
			return c.JSON(http.StatusInternalServerError, response)
		}

		return c.String(http.StatusOK, "Login Success")
	})

	token := <-okchan
	if token == "" {
		response := responses.NewResponse("error", "error while login", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	response := responses.GenerateResponses("ok", "login Success", []any{token})
	return c.JSON(http.StatusOK, response)
}
