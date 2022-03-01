package repository

import (
	"fmt"
	"github.com/OMUHA/oauwebscrapper/app/models/necta"
	"github.com/OMUHA/oauwebscrapper/config"
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
		return  false
	}
	return  school.ID  > 0
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


func  StoreSchoolToDB(school necta.School) {
	db := config.GetDBInstance()
	err := db.Model(&necta.School{}).Create(&school)
	if err.Error != nil {
		fmt.Println("Error while creating school", err.Error)
	}
}

func SearchSchoolDB(centerNo string) (necta.School,error) {
	db := config.GetDBInstance()
	school := necta.School{}
	err := db.Model(&necta.School{}).Where("center_no = ?", centerNo).First(&school).Error
	return school,err
}