package config

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"shopperia/src/API"
	commonMiddlewares "shopperia/src/common/middlewares"
)

func NewHttp(db *sql.DB) *echo.Echo {
	e := echo.New()

	mailConfig := NewEmailConfig()

	account, name, password, host, serverName := mailConfig.GetFields()

	API.BuildAPI(e, db, account, name, password, host, serverName)

	// e.Use(middleware.Recover())

	e.Use(commonMiddlewares.Logger)

	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Use(middleware.CORSWithConfig(corsConfig))

	return e
}
