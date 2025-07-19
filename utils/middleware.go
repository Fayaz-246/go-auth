package utils

import (
	"github.com/gofiber/fiber/v3"
)

func extractUserClaims(ctx fiber.Ctx, jwtKey string) (map[string]interface{}, error) {
	tokenStr := ctx.Get("AuthToken")
	if tokenStr == "" {
		return nil, fiber.ErrUnauthorized
	}

	claims, err := Verify(tokenStr, jwtKey)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	return claims, nil
}

func JWTMiddleware(jwtKey string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		claims, err := extractUserClaims(ctx, jwtKey)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or missing AuthToken",
			})
		}

		ctx.Locals("user", claims)
		return ctx.Next()
	}
}

func RequireRole(role string, jwtKey string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		claims, err := extractUserClaims(ctx, jwtKey)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or missing AuthToken",
			})
		}

		userRole, ok := claims["role"].(string)
		if !ok || userRole != role {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		ctx.Locals("user", claims)
		return ctx.Next()
	}
}

