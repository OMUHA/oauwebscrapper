package necta

import (
	"gorm.io/gorm"
)

type School struct {
	gorm.Model
	CenterNo string `json:"center_no" gorm:"type:varchar(255);uniqueIndex:candidateNumber_index;not null"`
	Name string `json:"name" gorm:"type:varchar(255);"`
	Region string `json:"region" gorm:"type:varchar(255);"`
}


type StudentResult struct {
	gorm.Model
	ID int `json:"id"`
	CandidateNo string `json:"candidate_no" gorm:"type:varchar(10);not null;uniqueIndex:candidateNumber_index"`
	CandidateFName string `json:"candidate_fname" gorm:"type:varchar(255); null"`
	CandidateMName string `json:"candidate_mname" gorm:"type:varchar(255); null"`
	CandidateLName string `json:"candidate_lname" gorm:"type:varchar(255); null"`
	CandidateEmail string `json:"candidate_email"  gorm:"type:varchar(255); null"`
	CandidatePhone string `json:"candidate_phone"  gorm:"type:varchar(255); null"`
	CandidateAddress string `json:"candidate_address"  gorm:"type:varchar(255); null"`
	CandidateGender string `json:"candidate_gender"  gorm:"type:varchar(3); null"`
	ResultsRaw string `json:"results_raw"`
	AggregatePoints string `json:"aggregate_points"  gorm:"type:varchar(10);not null"`
	ResultDivision string `json:"result_division"  gorm:"type:varchar(6);not null"`
	SchoolName string `json:"SchoolName"  gorm:"type:varchar(155); null"`
	CenterNo string `json:"center_no"  gorm:"type:varchar(10);not null"`
	ExamType string `json:"exam_type" gorm:"type:varchar(10);not null"`
	ExamYear string `json:"exam_year" gorm:"type:varchar(10);not null;uniqueIndex:candidateNumber_index"`
}
