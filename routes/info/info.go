package info

import (
	"goapi/utils"

	"github.com/gofiber/fiber/v3"
)

func RegRoutes(app *fiber.App, jwtkey string) {
	info := app.Group("info", utils.JWTMiddleware(jwtkey))

	info.Get("/me", func(ctx fiber.Ctx) error {
		user := ctx.Locals("user")

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "You are authorized",
			"user":    user,
		})
	})

}
