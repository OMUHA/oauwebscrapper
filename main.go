package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/OMUHA/oauwebscrapper/app/routers"
	"github.com/OMUHA/oauwebscrapper/config"
	"github.com/OMUHA/oauwebscrapper/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	config.InitDB()
	db := config.DB
	database.InitMigration(db)

	routers.Init(app)
	err = app.Listen(":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
		return
	}
}
