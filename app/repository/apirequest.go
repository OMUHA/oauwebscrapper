package repository

import (
	"crypto/tls"
	"encoding/xml"
	"github.com/OMUHA/oauwebscrapper/app/model"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func parseXML(xmlData string) (model.TCUResponseParameters, error) {
	var response model.TCUResponse

	err := xml.Unmarshal([]byte(xmlData), &response)
	if err != nil {
		return model.TCUResponseParameters{}, err
	}

	return response.ResponseParameters, nil
}

func GetCentersListing() ([]model.NectaSchool, error) {

	client := resty.New()
	var responResult struct {
		Response []model.NectaSchool          `json:"response"`
		Status   model.NectaApiResponseStatus `json:"status"`
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"api_key":"$2y$10$7BFbtDEWB2uac61b96WhlO7tAJp0p4bHbVYxhZgCe.D.WOGgHrG/2"}`).
		SetResult(&responResult).
		Post("https://api.necta.go.tz/api/secondary/centres")

	if err != nil {
		log.Fatal(err)
	}

	if resp.IsError() {
		log.Fatal(resp.RawResponse)
	}

	return responResult.Response, nil
}

func VerifyStudentAccount(indexNumber string) (model.TCUResponseParameters, error) {
	var request = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
		"<Request>" +
		"<UsernameToken><Username>MNZ</Username>" +
		"<SessionToken>1Hpn63x87qGSRTjr4OfE</SessionToken>" +
		"</UsernameToken><RequestParameters><f4indexno>" + indexNumber + "</f4indexno></RequestParameters></Request>"
	client := resty.New()
	var responResult string

	resp, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Content-Type", "application/xml").
		SetBody(request).
		SetResult(&responResult).
		Post("http://api.tcu.go.tz/applicants/checkStatus")
	if err != nil {
		log.Fatal(err)
	}

	if resp.IsError() {
		log.Println(resp.Request.Body)
		log.Println(resp.Request.URL)
		log.Fatal(resp.RawResponse)
	}

	response, err := parseXML(string(resp.Body()))

	return response, err

}

func UpdateStudentStatus(db *gorm.DB, students []model.TCUResponseParameters) {

	for _, student := range students {
		db.Where("f4index = ?", student.F4IndexNo).
			Updates(&model.ApplicantDetail{VerificationStatus: student.StatusDescription,
				VerificationCode: strconv.Itoa(student.StatusCode)})
	}

}

func GetStudentsListing(schoolNumber string) []model.NectaStudentDetail {

	client := resty.New()
	var responResult struct {
		Response []model.NectaStudentDetail   `json:"response"`
		Status   model.NectaApiResponseStatus `json:"status"`
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"centre_number":"` + schoolNumber + `","api_key":"$2y$10$7BFbtDEWB2uac61b96WhlO7tAJp0p4bHbVYxhZgCe.D.WOGgHrG/2"}`).
		SetResult(&responResult).
		Post("https://api.necta.go.tz/api/secondary/students")

	log.Println(resp.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.IsError() {
		log.Fatal(resp.RawBody())
	}

	return responResult.Response
}
