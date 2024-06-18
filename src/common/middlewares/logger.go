package commonMiddlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

func Logger(f echo.HandlerFunc) echo.HandlerFunc {
	timeSinceRequest := time.Now()
	return func(c echo.Context) error {
		requestOrigin := c.Request().Form
		path := c.Request().URL.Path
		status := c.Response().Status
		duration := time.Since(timeSinceRequest)

		loggerInfo := fmt.Sprintf("Request From: %v, To: %s, WithStatus: %s, serve in %s", requestOrigin, path, status, duration)

		log.Print(loggerInfo)

		return f(c)
	}
}
