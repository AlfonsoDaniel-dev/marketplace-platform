package userHandlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"shopperia/src/common/models"
	"shopperia/src/common/responses"
	"strconv"
	"time"
)

var loginStatus string

func (H *Handler) UserLogin(c echo.Context) error {
	form := models.Login{}

	err := c.Bind(&form)
	if err != nil {
		response := responses.NewResponse("error", "Bad structured Request", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	loginProgress := H.Service.LoginUser(form)
	if loginProgress.Error != nil {
		response := responses.NewResponse("error", "error while login", nil)
		return c.JSON(http.StatusBadRequest, response)
	} else if !loginProgress.HasTwoStepsVerification {
		response := responses.NewResponse("ok", "login success", loginProgress.Token)

		return c.JSON(http.StatusOK, response)
	}

	e := c.Echo()

	limit, err := strconv.Atoi(os.Getenv("MAX_LOGIN_CONFIRMATION_TIME"))
	if err != nil {
		limit = 30
	}

	confirmationChan := make(chan bool)

	go e.GET("api/user/login/confirm/:email/:token", func(c echo.Context) error {
		token := c.Param("token")

		token, err := H.Service.HandleTSVConfirmation(form, token)
		if err != nil {
			loginProgress.Error = err
			response := responses.NewResponse("error", "failed to login", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}

		confirmationChan <- true
		loginProgress.Token = token

		response := responses.NewResponse("ok", "Login success", nil)
		return c.JSON(http.StatusOK, response)
	})

	maxTime := time.Duration(limit) * time.Minute

	select {
	case <-confirmationChan:
		if loginProgress.Error != nil {
			response := responses.NewResponse("error", "login failed", err)
			return c.JSON(http.StatusInternalServerError, response)
		}
		response := responses.NewResponse("ok", "login success", loginProgress.Token)
		return c.JSON(http.StatusOK, response)
	case <-time.After(maxTime):
		response := responses.NewResponse("error", "login timeout no confirmation", nil)
		return c.JSON(http.StatusRequestTimeout, response)
	}

	response := responses.GenerateResponses("ok", "login Success", []any{loginProgress.Token})
	return c.JSON(http.StatusOK, response)
}
