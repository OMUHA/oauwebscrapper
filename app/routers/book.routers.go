package routers

import (
	"github.com/OMUHA/oauwebscrapper/app/controllers"
	"github.com/OMUHA/oauwebscrapper/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Book(app *fiber.App) {
	user := app.Group("/books")

	user.Get("/", middlewares.ExampleMiddleware, controllers.FetchAllBooks) // contoh menggunakan middleware
	user.Post("/", controllers.CreateBook)
}
