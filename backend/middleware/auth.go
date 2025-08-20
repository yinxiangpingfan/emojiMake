package middleware

import (
	"emoji-maker-backend/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected 是一个认证中间件
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwtService := services.NewJWTService()
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(services.APIResponse{
				Code:    1,
				Message: "Missing or malformed JWT",
			})
		}

		// Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(services.APIResponse{
				Code:    1,
				Message: "Missing or malformed JWT",
			})
		}
		tokenString := parts[1]

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(services.APIResponse{
				Code:    1,
				Message: "Invalid or expired JWT",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(services.APIResponse{
				Code:    1,
				Message: "Invalid JWT claims",
			})
		}

		userID := int64(claims["id"].(float64))
		c.Locals("userID", userID)
		return c.Next()
	}
}
