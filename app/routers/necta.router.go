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

func SinkingAPI(appRoutes *fiber.App) {
	appRoutes.Get("/download", controllers.DownloadAppData)
	appRoutes.Get("/verify_list", controllers.VerifyStudentList)
}
