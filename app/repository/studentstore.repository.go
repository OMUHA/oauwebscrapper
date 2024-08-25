package repository

import (
	"fmt"
	"github.com/OMUHA/oauwebscrapper/app/model"
	"github.com/OMUHA/oauwebscrapper/app/models/necta"
	"github.com/OMUHA/oauwebscrapper/config"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
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

func ExtractGrades(input string) map[string]string {
	grades := make(map[string]string)

	lines := strings.Split(input, "   ")

	for _, line := range lines {
		parts := strings.Split(line, " - ")
		if len(parts) == 2 {
			subject := strings.TrimSpace(parts[0])
			grade := strings.Trim(parts[1], "' ")
			grades[strings.ToUpper(subject)] = strings.ToUpper(grade)
		}
	}

	return grades
}

func CalculateGradePoint(grade string) int {
	switch grade {
	case "A":
		return 5
	case "B":
		return 4
	case "C":
		return 3
	case "D":
		return 2
	case "E":
		return 1
	case "F":
		return 0
	default:
		return -1
	}

}

func UpdateStudentResults(db *gorm.DB, year int) {
	var students []necta.StudentResult
	log.Printf(" starting %v ", year)

	err := db.Model(&necta.StudentResult{}).Where("exam_year = ?", strconv.Itoa(year)).Find(&students).Error
	if err != nil {
		log.Printf("Error while Fetching students %v \n", err)
	} else {
		var studentsUpdated []necta.StudentResult
		for _, student := range students {
			grades := ExtractGrades(student.ResultsRaw)
			for subject, grade := range grades {
				switch subject {
				case "B/MATH":
					student.Bmath = grade
					student.BmathPts = CalculateGradePoint(grade)
					break
				case "ENGL":
					student.Eng = grade
					student.EngPts = CalculateGradePoint(grade)
					break
				case "BIO":
					student.Bio = grade
					student.BioPts = CalculateGradePoint(grade)
					break
				case "PHY":
					student.Phy = grade
					student.PhyPts = CalculateGradePoint(grade)
					break
				case "CHEM":
					student.Chem = grade
					student.ChemPts = CalculateGradePoint(grade)
				}
			}
			if strings.Contains(strings.ToUpper(student.CandidateNo), "S") {
				student.CandidateType = "S"
			} else {
				student.CandidateType = "P"
			}

			var err3 = db.Model(&necta.StudentResult{}).Where("id = ?", student.ID).Updates(&student).Error
			if err3 != nil {
				log.Printf("Error updating student %s: %v \n", student.ID, student)
			} else {
				log.Printf("Student updated successfully %s \n", student.IndexNo)
			}
			//studentsUpdated = append(studentsUpdated, student)

		}
		log.Printf("Processed students %d total \n", len(studentsUpdated))
		/*if len(studentsUpdated) > 0 {
			updateNectaStudent(db, studentsUpdated)
		}*/
	}
}

func updateNectaStudent(db *gorm.DB, students []necta.StudentResult) {
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, std := range students {
			var err2 = db.Model(&necta.StudentResult{}).Where("id = ?", std.ID).Updates(&std).Error
			if err2 != nil {
				log.Printf("errors %v", err2)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error updating student: %v \n", err)
	} else {
		log.Printf("Updated student total : %d \n", len(students))
	}
	panic(len(students))
}

func CountTotalNectaStudent(db *gorm.DB) int64 {

	var totalCount int64
	db.Model(&necta.StudentResult{}).Count(&totalCount)
	return totalCount
}

func CheckNectaSchoolHasStudents(centerNo string, year string, exam_type string) bool {
	db := config.GetDBInstance()
	school := necta.StudentResult{}
	var totalStudents int64
	db.Model(&school).Where("center_no = ? and exam_year = ? and exam_type = ?", centerNo, year, exam_type).Count(&totalStudents)
	return totalStudents > 0
}

func CheckSchoolHasStudents(db *gorm.DB, school model.NectaSchool) bool {
	var countedStudents int64
	db.Model(&model.NectaStudentDetail{}).Where("center_number = ?", school.Number).Count(&countedStudents)
	return countedStudents > 0
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

func CreateNectaSchoolStudents(db *gorm.DB, students []model.NectaStudentDetail, centerNo string, centerID uint, regYear string) {

	var studentsCreate []model.NectaStudentDetail
	for _, student := range students {
		var existed model.NectaStudentDetail
		db.Model(&model.NectaStudentDetail{}).Where("psle_number = ? and center_number = ?", student.PsleNumber, centerNo).First(&existed)
		if existed.ID > 0 {
			log.Printf("Student exists %s \n", student.PsleNumber)
		} else {
			student.CenterNumber = centerNo
			student.CenterId = centerID
			student.RegYear = regYear
			student.Disabilities = nil
			student.Difficulties = nil
			studentsCreate = append(studentsCreate, student)
		}
	}
	db.Model(&model.NectaStudentDetail{}).Create(&studentsCreate)
}

func CreateStudentDetails(db *gorm.DB, student model.ApplicantDetail) {
	var err = db.Model(&model.ApplicantDetail{}).Create(&student).Error
	if err != nil {
		log.Fatalf(" Error saving student %s ", err)
	}
	log.Printf("Exttacted %s - %s - %s \n", student.HliID, student.F4index, student.F6Index)
}

func GetApplicantData(db *gorm.DB) []model.ApplicantDetail {
	var students []model.ApplicantDetail
	db.Model(&model.ApplicantDetail{}).Find(&students)
	return students
}

func GetTotalStuentDetaisl(db *gorm.DB) int64 {
	var total int64
	db.Model(&model.ApplicantDetail{}).Count(&total)
	return total
}

func GetApplicantDataLimited(db *gorm.DB, start, limit int) []model.ApplicantDetail {
	var students []model.ApplicantDetail
	err := db.Model(&model.ApplicantDetail{}).Offset(start).Limit(limit).Find(&students).Error

	if err != nil {
		log.Printf(" errors %s", err.Error())
	}
	log.Printf(" total STudents %d", len(students))
	return students
}
