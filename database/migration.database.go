package database

import (
	"github.com/OMUHA/oauwebscrapper/app/model"
	"github.com/OMUHA/oauwebscrapper/app/models"
	"github.com/OMUHA/oauwebscrapper/app/models/necta"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	err := db.AutoMigrate(&model.ApplicantDetail{}, &model.FilteredApplicantDetail{}, &model.NectaSchool{}, &model.NectaStudentDetail{}, &models.User{}, &necta.StudentResult{}, &necta.School{})
	if err != nil {
		//panic(err)
	}
}
