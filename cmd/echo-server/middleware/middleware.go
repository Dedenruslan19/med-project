package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func forbiddenResponse(c echo.Context) error {
	return c.JSON(http.StatusForbidden, map[string]interface{}{"message": http.StatusText(http.StatusForbidden)})
}

func JWTMiddleware(jwtSign string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if strings.Contains(c.Request().URL.Path, "/login") {
				return next(c)
			}

			signature := strings.Split(c.Request().Header.Get("Authorization"), " ")
			if len(signature) < 2 {
				return forbiddenResponse(c)
			}
			if signature[0] != "Bearer" {
				return forbiddenResponse(c)
			}

			claim := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(signature[1], claim, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(jwtSign), nil
			})
			if err != nil {
				return forbiddenResponse(c)
			}

			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwt.SigningMethodHS256 {
				return forbiddenResponse(c)
			}

			expAt, err := claim.GetExpirationTime()
			if err != nil {
				return forbiddenResponse(c)
			}

			if time.Now().After(expAt.Time) {
				return forbiddenResponse(c)
			}

			userID, _ := claim["id"].(string)
			role, _ := claim["role"].(string)
			c.Set("id", userID)
			c.Set("role", role)

			return next(c)
		}
	}
}

func ACLMiddleware(rolesMap map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get("role").(string)
			if role == "" {
				return next(c)
			}

			if rolesMap[role] {
				return next(c)
			}

			return forbiddenResponse(c)
		}
	}
}

func JwtEchoMiddleware(jwtSign string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSign),
	})
}

func ValidateContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "invalid content type"})
		}
		return next(c)
	}
}

// GetUserID returns id claim as int64 (handles string/number stored in context)
func GetUserID(c echo.Context) (int64, bool) {
	v := c.Get("id")
	if v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case string:
		if t == "" {
			return 0, false
		}
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, false
		}
		return id, true
	case float64:
		return int64(t), true
	case int64:
		return t, true
	case int:
		return int64(t), true
	default:
		return 0, false
	}
}

// GetRole returns the role string stored in context
func GetRole(c echo.Context) (string, bool) {
	v := c.Get("role")
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
