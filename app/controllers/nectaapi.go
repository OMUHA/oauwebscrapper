package controllers

import (
	"log"

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
			updatedSchool := repository.CreateNectaSchool(db, school, i)

			hasStudents := repository.CheckSchoolHasStudents(db,school);

			if hasStudents {
				log.Println("school has students downloaded $s",school.Number)
			}else{

			if updatedSchool.ID > 0 {
				students := repository.GetStudentsListing(school.Number)
				if len(students) > 0 {
					repository.CreateNectaSchoolStudents(db, students, school.Number, updatedSchool.ID)
				} else {
					log.Println("no students found on %s", school.Number)
				}
			}
			}
		}
	}
	return nil
}
