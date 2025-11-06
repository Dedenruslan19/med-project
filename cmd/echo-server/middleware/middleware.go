package middleware

import (
	"net/http"
	"strings"

	"Dedenruslan19/med-project/util"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid Authorization header format",
			})
		}

		tokenStr := parts[1]
		claims, err := util.ValidateJWT(tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid or expired token",
			})
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		return next(c)
	}
}

func ValidateContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid content type",
			})
		}
		return next(c)
	}
}
