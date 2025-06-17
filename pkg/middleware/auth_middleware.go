package middleware

import (
	"context"
	"strings"

	service "booking/internal/user"
	"booking/pkg/jwt"
	"booking/pkg/response"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(userService service.UserServiceInterface, jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, "authorization header is required", nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Unauthorized(c, "invalid authorization header format", nil)
			}

			tokenString := parts[1]
			claims, err := jwt.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				return response.Unauthorized(c, "invalid token", err)
			}

			// Periksa token di Redis
			ctx := context.Background()
			storedToken, err := userService.GetStoredToken(ctx, claims.UserID)
			if err != nil {
				// Jika token tidak ditemukan di Redis, berarti user sudah logout
				return response.Unauthorized(c, "token has been revoked or expired", nil)
			}

			if storedToken != tokenString {
				// Jika token berbeda, berarti user menggunakan token yang tidak valid
				return response.Unauthorized(c, "token is invalid or has been revoked", nil)
			}

			user, err := userService.GetUserByID(ctx, claims.UserID)
			if err != nil {
				return response.Unauthorized(c, "user not found", err)
			}

			if user == nil {
				return response.Unauthorized(c, "user not found", nil)
			}

			// Pastikan user adalah pointer yang valid sebelum disimpan ke context
			if user != nil {
				c.Set("user", user)
				c.Set("user_id", claims.UserID)
			} else {
				return response.Unauthorized(c, "invalid user data", nil)
			}

			return next(c)
		}
	}
}
