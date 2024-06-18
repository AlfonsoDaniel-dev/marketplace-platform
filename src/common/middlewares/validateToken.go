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
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Not authorized to access this resource"})
		}

		return f(c)
	}
}
