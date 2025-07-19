package admin

import (
	"goapi/utils"

	"github.com/gofiber/fiber/v3"
)

func RegRoutes(app *fiber.App, jwtkey string) {
	admin := app.Group("admin", utils.RequireRole("admin", jwtkey))

	admin.Get("/me", func(ctx fiber.Ctx) error {
		user := ctx.Locals("user")

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "You are an admin",
			"user":    user,
		})

	})
}
