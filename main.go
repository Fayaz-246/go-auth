package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"goapi/routes"
	"goapi/structs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtkey := os.Getenv("JWT_KEY")

	db.AutoMigrate(&structs.User{})

	routes.Setup(app, db, jwtkey)
	log.Fatal(app.Listen(":3000"))
}
