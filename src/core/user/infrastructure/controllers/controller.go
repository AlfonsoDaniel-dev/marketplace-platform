package Usercontroller

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	commonMiddlewares "shopperia/src/common/middlewares"
	UserServices "shopperia/src/core/user/app/Services"
	userDomain "shopperia/src/core/user/domain"
	User_email "shopperia/src/core/user/infrastructure/Email"
	userHandlers "shopperia/src/core/user/infrastructure/controllers/handlers"
	Userstorage "shopperia/src/core/user/infrastructure/storage"
)

type Handler interface {
	HelloWorld(c echo.Context) error
	Register(c echo.Context) error
	UserLogin(c echo.Context) error
	UpdateUserName(c echo.Context) error
	UpdateUserFirstName(c echo.Context) error
	UpdateUserLastName(c echo.Context) error
	UpdateUserEmail(c echo.Context) error
	UpdateUserPassword(c echo.Context) error
	UpdateTsvStatus(c echo.Context) error
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
	c.userPublicRoutes()
	c.userPrivateRoutes()
}

func (c *UserController) testRoutes() {
	tester := c.echo.Group("api/v1/test")

	tester.GET("/helloWorld", c.Handler.HelloWorld)
}

func (c *UserController) userPublicRoutes() {
	public := c.echo.Group("/api/v1/user")

	public.POST("/register", c.Handler.Register)
	public.POST("/login", c.Handler.UserLogin)
}

func (c *UserController) userPrivateRoutes() {
	private := c.echo.Group("/api/v1/user")
	private.Use(commonMiddlewares.ValidateToken)

	// config routes
	config := private.Group("/config")
	config.PUT("/update/user_name", c.Handler.UpdateUserName)
	config.PUT("/update/first_name", c.Handler.UpdateUserFirstName)
	config.PUT("/update/last_name", c.Handler.UpdateUserLastName)
	config.PUT("/update/email", c.Handler.UpdateUserEmail)
	config.PUT("/update/password", c.Handler.UpdateUserPassword)
	config.PUT("/update/tsv", c.Handler.UpdateTsvStatus)
}
