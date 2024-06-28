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

func (H *Handler) UpdateUserFirstName(c echo.Context) error {
	form := UserDTO.UpdateFirstName{}

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

	err = H.Service.UpdateUserFirstName(form)
	if err != nil {
		fmt.Println(form.FirstName)
		response := responses.NewResponse("error", "update user name failed", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := responses.GenerateResponses("ok", "first name updated successfully", nil)
	return c.JSON(http.StatusOK, response)
}

func (H *Handler) UpdateUserLastName(c echo.Context) error {
	form := UserDTO.UpdateLastName{}

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

	err = H.Service.UpdateUserLastName(form)
	if err != nil {
		fmt.Println(err)
		response := responses.NewResponse("error", "update user name failed", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := responses.GenerateResponses("ok", "last name updated successfully", nil)
	return c.JSON(http.StatusOK, response)
}

func (H *Handler) UpdateUserEmail(c echo.Context) error {
	form := UserDTO.UpdateEmail{}

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
	err = H.Service.UpdateUserEmail(form)
	if err != nil {
		response := responses.NewResponse("error", "update user email failed", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := responses.GenerateResponses("ok", "email updated successfully", nil)
	return c.JSON(http.StatusOK, response)
}

func (H *Handler) UpdateUserPassword(c echo.Context) error {
	form := UserDTO.UpdatePassword{}

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

	_, err = H.Service.SendPasswordEmailConfirmation(email)
	if err != nil {
		response := responses.NewResponse("error", "send password email confirmation failed", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	e := c.Echo()

	okChan := make(chan bool)

	go e.GET("/api/v1/user/config/update/password/confirm/:email/:token", func(c echo.Context) error {
		email := c.Param("email")
		token := c.Param("token")

		ok, err := H.Service.CheckToken(email, token)
		if err != nil {
			okChan <- false
			response := responses.NewResponse("error", "error while checking token", nil)
			return c.JSON(http.StatusInternalServerError, response)
		} else if !ok {
			okChan <- false
			response := responses.NewResponse("error", "invalid token", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		err = H.Service.CleanToken(email)
		if err != nil {
			okChan <- false
			response := responses.NewResponse("error", "error while cleaning token", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}

		okChan <- true

		response := responses.NewResponse("ok", "password updated successfully continue on your session", nil)
		return c.JSON(http.StatusOK, response)
	})

	confirm := <-okChan

	if !confirm {
		response := responses.NewResponse("error", "password change not authorized", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	form.Common.Email = email

	err = H.Service.UpdateUserPassword(form)
	if err != nil {
		fmt.Println(err)
		response := responses.NewResponse("error", "update user password failed", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := responses.GenerateResponses("ok", "password updated successfully", nil)
	return c.JSON(http.StatusOK, response)
}

func (H *Handler) UpdateTsvStatus(c echo.Context) error {
	form := UserDTO.UpdateTsvStatus{}

	err := c.Bind(&form)
	if err != nil {
		fmt.Println(err)
	}

	token := c.Request().Header.Get("authorization")
	email, err := auth.GetFromToken(token, "email")
	if err != nil {
		response := responses.NewResponse("error", "invalid token", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	form.Common.Email = email

	okChan := make(chan bool)

	e := c.Echo()

	go e.GET("/api/user/config/tsv/confirm/:email/:token", func(c echo.Context) error {
		email := c.Param("email")
		token := c.Param("token")

		ok, err := H.Service.CheckToken(email, token)
		if err != nil {
			okChan <- false
			response := responses.NewResponse("error", "error while checking token", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}

		if !ok {
			okChan <- false
			response := responses.NewResponse("error", "invalid token", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		okChan <- true

		response := responses.GenerateResponses("ok", "tsv status updated successfully", nil)
		return c.JSON(http.StatusOK, response)
	})

	ok := <-okChan
	if !ok {
		response := responses.NewResponse("error", "Cannot Update Tsv status", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	err = H.Service.UpdateUserTsvConfig(form)
	if err != nil {
		fmt.Println(err)
	}

	response := responses.GenerateResponses("ok", "tsv status updated successfully", nil)
	return c.JSON(http.StatusOK, response)
}
