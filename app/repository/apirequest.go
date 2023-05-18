package repository

import (
	"github.com/OMUHA/oauwebscrapper/app/model"
	"github.com/go-resty/resty/v2"
	"log"
)

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
		log.Fatal(resp.RawResponse)
	}

	return responResult.Response
}
