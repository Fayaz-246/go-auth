package utils

import	"github.com/gofiber/fiber/v3"

func JWTMiddleware(jwtkey string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		tokenStr := ctx.Get("AuthToken")
		if tokenStr == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "AuthToken header missing",
			})
		}

		claims, err := Verify(tokenStr, jwtkey)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		ctx.Locals("user", claims)

		return ctx.Next()
	}
}

func RequireRole(role string, jwtkey string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		tokenStr := ctx.Get("AuthToken")
		if tokenStr == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "AuthToken header missing",
			})
		}

		claims, err := Verify(tokenStr, jwtkey)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		ctx.Locals("user", claims)

		if claims["role"] != role {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}
		return ctx.Next()
	}
}
