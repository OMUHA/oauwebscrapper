package controllers

import (
	"github.com/OMUHA/oauwebscrapper/app/models"
	"log"
	"net/http"

	"github.com/OMUHA/oauwebscrapper/app/repository"
	"github.com/OMUHA/oauwebscrapper/config"
	"github.com/gofiber/fiber/v2"
)

func GetNectaSchoolListing(ctx *fiber.Ctx) error {
	listing, err := repository.GetCentersListing()
	if err != nil {
		return err
	}
	db := config.GetDBInstance()

	if len(listing) > 0 {
		for i, school := range listing {
			repository.CreateNectaSchool(db, school, i)
			// updatedSchool := repository.CreateNectaSchool(db, school, i)
			/*hasStudents := repository.CheckSchoolHasStudents(db, school)
			if hasStudents {
				log.Println("school has students downloaded $s", school.Number)
			} else {
				if updatedSchool.ID > 0 {
					students := repository.GetStudentsListing(school.Number)
					if len(students) > 0 {
						repository.CreateNectaSchoolStudents(db, students, school.Number, updatedSchool.ID)
					} else {
						log.Println("no students found on %s", school.Number)
					}
				}
			}*/
		}
	}
	return nil
}
func GetNectaStudentListing(ctx *fiber.Ctx) error {
	schoolsList, err := repository.FindAllNectaSchools()
	var yearID = ctx.Params("yearID")
	var response models.Response
	if err != nil {
		response.Message = err.Error()
		response.Status = http.StatusInternalServerError
		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}
	db := config.GetDBInstance()
	for _, school := range schoolsList {
		hasStudents := repository.CheckSchoolHasStudents(db, school)
		if hasStudents {
			log.Printf("school has students downloaded %s\n", school.Number)
		} else {
			if school.ID > 0 {
				students := repository.GetStudentsListing(school.Number)
				if len(students) > 0 {
					repository.CreateNectaSchoolStudents(db, students, school.Number, school.ID, yearID)
				} else {
					log.Printf("no students found on %s\n", school.Number)
				}
			}
		}
	}
	response.Message = "Response data"
	response.Data = schoolsList
	response.Status = 200
	return ctx.Status(200).JSON(response)
}
