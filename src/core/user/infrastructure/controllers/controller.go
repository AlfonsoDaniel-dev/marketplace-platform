package Usercontroller

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	UserServices "shopperia/src/core/user/app/Services"
	userDomain "shopperia/src/core/user/domain"
	User_email "shopperia/src/core/user/infrastructure/Email"
	userHandlers "shopperia/src/core/user/infrastructure/controllers/handlers"
	Userstorage "shopperia/src/core/user/infrastructure/storage"
)

type Handler interface {
	HelloWorld(c echo.Context) error
	Register(c echo.Context) error
}

type UserController struct {
	Handler
	echo *echo.Echo
}

func NewController(e *echo.Echo, db *sql.DB, accountEmail, accountName, password, host, serverName string) *UserController {
	psqlUser := Userstorage.NewPsqlUser(db)
	emailSender := User_email.NewEmailSender(accountEmail, accountName, password, host, serverName)
	UserDomain := userDomain.NewUserDomain(psqlUser, emailSender)
	Services := UserServices.NewService(UserDomain)

	handler := userHandlers.NewHandler(Services)

	return &UserController{
		Handler: handler,
		echo:    e,
	}
}

func (c UserController) BuildRoutes() {

	c.testRoutes()
	c.UserPublicRoutes()
}

func (c *UserController) testRoutes() {
	tester := c.echo.Group("/test")

	tester.GET("/helloWorld", c.Handler.HelloWorld)
}

func (c *UserController) UserPublicRoutes() {
	public := c.echo.Group("/api/user")

	public.POST("/register", c.Handler.Register)
}
