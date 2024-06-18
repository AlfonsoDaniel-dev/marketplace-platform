package API

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	Usercontroller "shopperia/src/core/user/infrastructure/controllers"
)

func BuildAPI(e *echo.Echo, db *sql.DB, accountEmail string, accountName string, password string, host string, serverName string) {
	userController := Usercontroller.NewController(e, db, accountEmail, accountName, password, host, serverName)

	userController.BuildRoutes()
}
