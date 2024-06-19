package userHandlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"shopperia/src/common/models"
	"shopperia/src/common/responses"
)

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

	link, err := H.Service.SendLoginConfirmation(form.Email)
	if err != nil {
		response := responses.NewResponse("error", "cannot send email verification", err)
		return c.JSON(http.StatusBadRequest, response)
	}

	e := c.Echo()

	okChan := make(chan bool)

	go e.GET(link, func(c echo.Context) error {
		okChan <- true
		return c.String(http.StatusOK, "Response get, ")
	})

	LoginTsvStatus := <-okChan

	if !LoginTsvStatus {
		response := responses.NewResponse("error", "login cannot be complete", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	token, err := H.Service.LoginUser(form)
	if err != nil {
		response := responses.NewResponse("error", "error while login", err)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := responses.GenerateResponses("ok", "login Success", []any{token})
	return c.JSON(http.StatusOK, response)
}
