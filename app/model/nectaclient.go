package model

import "gorm.io/gorm"

type NectaApiResponseStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type NectaApiResponse struct {
	Response []interface{}          `json:"response"`
	Status   NectaApiResponseStatus `json:"status"`
}

type NectaSchool struct {
	gorm.Model
	Number             string `json:"number" gorm:"size:200"`
	RegistrationNumber string `json:"registration_number" gorm:"size:200"`
	Name               string `json:"name" gorm:"size:200"`
	DistrictId         int    `json:"district_id"`
	WardId             int    `json:"ward_id"`
	Ownership          string `json:"ownership"  gorm:"size:200"`
	PrincipalName      string `json:"principal_name"  gorm:"size:200"`
	PrincipalPhone     string `json:"principal_phone"  gorm:"size:200"`
	ContactOne         string `json:"contact_one"  gorm:"size:200"`
	ContactTwo         string `json:"contact_two"  gorm:"size:200"`
	DistrictDistance   string `json:"district_distance"  gorm:"size:200"`
	CenterType         string `json:"center_type"  gorm:"size:200"`
	PostalAddress      string `json:"postal_address"  gorm:"size:200"`
	IsGovernment       int    `json:"is_government"`
}

type NectaStudentDetail struct {
	gorm.Model
	ID                      int           `json:"id"`
	PremNumber              string        `json:"prem_number"  gorm:"size:200"`
	FirstName               string        `json:"first_name"   gorm:"size:200"`
	OtherNames              string        `json:"other_names"   gorm:"size:200"`
	Surname                 string        `json:"surname"   gorm:"size:200"`
	DateOfBirth             string        `json:"date_of_birth"   gorm:"size:200"`
	Sex                     string        `json:"sex" gorm:"size:4"`
	PsleNumber              string        `json:"psle_number"   gorm:"size:200"`
	PsleYear                string        `json:"psle_year"   gorm:"size:200"`
	IDNumber                string        `json:"id_number"   gorm:"size:200"`
	Photo                   string        `json:"photo"   gorm:"size:200"`
	RegistrationType        int           `json:"registration_type"`
	IsRepeater              int           `json:"is_repeater"`
	Status                  int           `json:"status"`
	IsTanzanian             int           `json:"is_tanzanian"`
	IsOrphan                int           `json:"is_orphan"`
	RegistrationDate        string        `json:"registration_date"   gorm:"size:200"`
	BirthCertificateNo      string        `json:"birth_certificate_no"   gorm:"size:200"`
	DistanceFromHome        string        `json:"distance_from_home"   gorm:"size:200"`
	AddressLine1            string        `json:"address_line_1"   gorm:"size:200"`
	AddressLine2            string        `json:"address_line_2"   gorm:"size:200"`
	AddressLine3            string        `json:"address_line_3"   gorm:"size:200"`
	DateOfExit              string        `json:"date_of_exit"   gorm:"size:200"`
	FormID                  int           `json:"form_id"`
	ClassCode               string        `json:"class_code"   gorm:"size:200"`
	Citizenship             int           `json:"citizenship"   gorm:"size:200"`
	PhysicalAddress         string        `json:"physical_address"   gorm:"size:200"`
	Domicile                string        `json:"domicile"   gorm:"size:200"`
	IsAdmitted              int           `json:"is_admitted"`
	GuardianNIN             string        `json:"guardian_NIN"   gorm:"size:200"`
	GuardianSex             string        `json:"guardian_sex"     gorm:"size:200"`
	GuardianPhone           string        `json:"guardian_phone"   gorm:"size:200"`
	GuardianName            string        `json:"guardian_name"   gorm:"size:200"`
	GuardianRelation        string        `json:"guardian_relation"   gorm:"size:200"`
	GuardianAddress         string        `json:"guardian_address"   gorm:"size:200"`
	GuardianEmail           string        `json:"guardian_email"   gorm:"size:200"`
	GuardianOccupation      string        `json:"guardian_occupation"   gorm:"size:200"`
	GuardianPhysicalAddress string        `json:"guardian_physical_address"   gorm:"size:200"`
	ParentingStatus         string        `json:"parenting_status"`
	Difficulties            []interface{} `json:"difficulties" gorm:"-" `
	Disabilities            []interface{} `json:"disabilities"  gorm:"-" `
	CenterNumber            string        `json:"center_number" gorm:"size:10"`
	CenterId                uint          `json:"center_id"`
	RegYear                 string        `json:"reg_year" gorm:"size:10"`
}

type ApplicantDetail struct {
	gorm.Model
	HliID              string `json:"hli_id"`
	F4index            string `json:"f4_index"`
	F6Index            string `json:"f6_index"`
	AdmittedProgram    string `json:"admitted_program"`
	Fname              string `json:"fname"`
	Mname              string `json:"mname"`
	Lname              string `json:"lname"`
	Gender             string `json:"gender"`
	F4result           string `json:"f4result"`
	F6Result           string `json:"F6Result"`
	MobileNumber       string `json:"mobile_number"`
	EmailAddress       string `json:"email_address"`
	AdmissionStatus    string `json:"admission_status"`
	Programs           string `json:"programs"`
	Comment            string `json:"comment"`
	VerificationStatus string `json:"verification_status"`
	VerificationCode   string `json:"verification_code"`
}

type FilteredApplicantDetail struct {
	ApplicantDetail
}

type TCUResponseParameters struct {
	F4IndexNo         string `xml:"f4indexno"`
	StatusCode        int    `xml:"StatusCode"`
	StatusDescription string `xml:"StatusDescription"`
}

type TCUResponse struct {
	ResponseParameters []TCUResponseParameters `xml:"ResponseParameters"`
}
