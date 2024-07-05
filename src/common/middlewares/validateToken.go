package commonMiddlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"shopperia/src/auth"
)

func ValidateToken(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
		}

		_, err := auth.ValidateToken(tokenString)
		if err != nil {
			if err := c.Redirect(http.StatusUnauthorized, "/api/v1/user/login"); err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}
		}

		return f(c)
	}
}
