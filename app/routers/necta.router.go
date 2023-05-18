package routers

import (
	"github.com/OMUHA/oauwebscrapper/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func NectaScrapper(appRoutes *fiber.App) {
	appRoutes.Get("/necta/scrapper/:yearID", controllers.NectaCseeScrapper)
}

func NectaAPI(appRoutes *fiber.App) {
	appRoutes.Get("/necta/getSchools", controllers.GetNectaSchoolListing)
}
