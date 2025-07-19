package auth

import (
	"goapi/structs"
	"goapi/utils"

	"errors"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log"
)

func RegRoutes(app *fiber.App, db *gorm.DB, jwtkey string) {
	auth := app.Group("/auth")

	auth.Post("/signup", func(ctx fiber.Ctx) error {
		payload := new(structs.RegisterPayload)

		var existingUser structs.User

		err := ctx.Bind().Body(payload)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		err = db.Where("email = ? OR name = ?", payload.Email, payload.Name).First(&existingUser).Error

		if err == nil {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User already exists with same email or name.",
			})
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}

		validEmail := utils.CheckEmail(payload.Email)
		if !validEmail {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email format",
			})
		}

		user := structs.User{
			Email:        payload.Email,
			Name:         payload.Name,
			PasswordHash: hashedPassword,
			Role:         payload.Role,
		}

		dbUser := db.Create(&user)

		if dbUser.Error != nil {
			log.Println("DB insert failed:", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		token, err := utils.Sign(user, jwtkey)

		if err != nil {
			log.Fatal("Failed to sign JWT:", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to sign JWT"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
	})

	auth.Post("/login", func(ctx fiber.Ctx) error {
		payload := new(structs.LoginPayload)
		if err := ctx.Bind().Body(payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		var existingUser structs.User

		err := db.Where("email = ?", payload.Email).First(&existingUser).Error

		if err != nil {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User does not exist.",
			})
		}

		if !utils.CheckPasswordHash(payload.Password, existingUser.PasswordHash) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Invalid Password"})
		}

		token, err := utils.Sign(existingUser, jwtkey)

		if err != nil {
			log.Fatal("Failed to sign JWT:", err)
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
	})
}
