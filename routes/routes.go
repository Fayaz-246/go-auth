package routes

import (
	"goapi/routes/admin"
	"goapi/routes/auth"
	"goapi/routes/info"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, jwtkey string) {
	auth.RegRoutes(app, db, jwtkey)
	info.RegRoutes(app, jwtkey)
	admin.RegRoutes(app, jwtkey)
}
