package necta

import (
	"gorm.io/gorm"
)

type School struct {
	gorm.Model
	CenterNo string `json:"center_no" gorm:"type:varchar(255);uniqueIndex:candidateNumber_index;not null"`
	Name     string `json:"name" gorm:"type:varchar(255);"`
	Region   string `json:"region" gorm:"type:varchar(255);"`
}

type SchoolStatistic struct {
	gorm.Model
	CenterNo    string  `json:"center_no" gorm:"type:varchar(20)"`
	SubjectCode string  `json:"subject_code" gorm:"type:varchar(10)"`
	SubjectName string  `json:"subject_name" gorm:"type:varchar(60)"`
	Registered  int     `json:"registered" `
	Sat         int     `json:"sat"`
	Clean       int     `json:"clean"`
	Pass        int     `json:"pass"`
	GPA         float64 `json:"gpa"`
	RegRank     int     `json:"reg_rank"`
	NatRank     int     `json:"nat_rank"`
	NatRankTot  int     `json:"nat_tot"`
	RegRankTot  int     `json:"reg_rank_tot"`
	ExamType    string  `json:"exam_type" gorm:"type:varchar(10)"`
	ExamYear    int     `json:"exam_year"`
}

type StudentResult struct {
	gorm.Model
	ID               int    `json:"id"`
	IndexNo          string `json:"index_no" gorm:"type:varchar(50);uniqueIndex:candidateNumber_index"`
	CandidateNo      string `json:"candidate_no" gorm:"type:varchar(10);not null;index:candidateNumber_index"`
	CandidateFName   string `json:"candidate_fname" gorm:"type:varchar(255); null"`
	CandidateMName   string `json:"candidate_mname" gorm:"type:varchar(255); null"`
	CandidateLName   string `json:"candidate_lname" gorm:"type:varchar(255); null"`
	CandidateEmail   string `json:"candidate_email"  gorm:"type:varchar(255); null"`
	CandidatePhone   string `json:"candidate_phone"  gorm:"type:varchar(255); null"`
	CandidateAddress string `json:"candidate_address"  gorm:"type:varchar(255); null"`
	CandidateGender  string `json:"candidate_gender"  gorm:"type:varchar(3); null"`
	CandidateType    string `json:"candidate_type" gorm:"type:varchar(2); null"`
	ResultsRaw       string `json:"results_raw"`
	AggregatePoints  string `json:"aggregate_points"  gorm:"type:varchar(10);not null"`
	ResultDivision   string `json:"result_division"  gorm:"type:varchar(6);not null"`
	SchoolName       string `json:"SchoolName"  gorm:"type:varchar(155); null"`
	CenterNo         string `json:"center_no"  gorm:"type:varchar(10);not null"`
	ExamType         string `json:"exam_type" gorm:"type:varchar(10);not null"`
	ExamYear         string `json:"exam_year" gorm:"type:varchar(10);not null;index:candidateNumber_index"`
	Phy              string `json:"phy" gorm:"type:varchar(2)"`
	PhyPts           int    `json:"phy_pts"`
	Chem             string `json:"chem" gorm:"type:varchar(2)"`
	ChemPts          int    `json:"chem_pts"`
	Bio              string `json:"bio" gorm:"type:varchar(2)"`
	BioPts           int    `json:"bio_pts"`
	Bmath            string `json:"bmat" gorm:"type:varchar(2)"`
	Eng              string `json:"eng" gorm:"type:varchar(2)"`
	EngPts           int    `json:"eng_pts"`
	BmathPts         int    `json:"bmath"`
}
