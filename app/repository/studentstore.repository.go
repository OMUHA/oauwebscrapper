package repository

import (
	"fmt"
	"github.com/OMUHA/oauwebscrapper/app/model"
	"github.com/OMUHA/oauwebscrapper/app/models/necta"
	"github.com/OMUHA/oauwebscrapper/config"
	"gorm.io/gorm"
	"log"
)

func StoreStudentResults(studentResults []necta.StudentResult) error {
	studentResults = studentResults[1:]
	return StoreStudentResultsListToDB(studentResults)
}

func StoreSchool(school necta.School) error {
	StoreSchoolToDB(school)
	return nil
}

func CheckSchoolExists(centerNo string) bool {
	var school necta.School
	school, err := SearchSchoolDB(centerNo)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return school.ID > 0
}

func StoreStudentResultsListToDB(students []necta.StudentResult) error {
	db := config.GetDBInstance()
	err := db.Model(&necta.StudentResult{}).Create(&students).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func StoreSchoolToDB(school necta.School) {
	db := config.GetDBInstance()
	var schoolN necta.School
	db.Model(&necta.School{}).Where("center_no = ?", school.CenterNo).First(&schoolN)
	if schoolN.ID > 0 {
		fmt.Println("School existing")
	} else {
		err := db.Model(&necta.School{}).Create(&school)
		if err.Error != nil {
			fmt.Println("Error while creating school", err.Error)
		}
	}

}

func SearchSchoolDB(centerNo string) (necta.School, error) {
	db := config.GetDBInstance()
	school := necta.School{}
	err := db.Model(&necta.School{}).Where("center_no = ?", centerNo).First(&school).Error
	return school, err
}

func CreateNectaSchool(db *gorm.DB, schol model.NectaSchool, index int) model.NectaSchool {

	schoolN := model.NectaSchool{}

	db.Model(&model.NectaSchool{}).Where("number = ?", schol.Number).Find(&schoolN)
	if schoolN.ID > 0 {
		log.Println("school name exists %v", schol.Number)
		return schoolN
	} else {
		db.Model(&model.NectaSchool{}).Create(&schol)
	}
	return schol
}

func CreateNectaSchoolStudents(db *gorm.DB, students []model.NectaStudentDetail, centerNo string, centerID uint) {

	var studentsCreate []model.NectaStudentDetail
	for _, student := range students {
		var existed model.NectaStudentDetail
		db.Model(&model.NectaStudentDetail{}).Where("psle_number = ? and center_number = ?", student.PsleNumber, centerNo).First(&existed)
		if existed.ID > 0 {
			log.Println("Student exists %s ", student.PsleNumber)
		} else {
			student.CenterNumber = centerNo
			student.CenterId = centerID
			student.Disabilities = nil
			student.Difficulties = nil
			studentsCreate = append(studentsCreate, student)
		}

	}
	db.Model(&model.NectaStudentDetail{}).Create(&studentsCreate)
}
