package repository

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/OMUHA/oauwebscrapper/app/model"
	"github.com/OMUHA/oauwebscrapper/config"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

/*func parseXML(xmlData []byte) (model.TCUResponseParameters, error) {
	var response model.TCUResponse

	err := xml.Unmarshal(xmlData, &response)
	if err != nil {
		log.Fatal(err.Error())
		return model.TCUResponseParameters{}, err
	}

	return response.ResponseParameters, nil
}*/

func parseXML(xmlData string) (model.TCUResponseParameters, error) {
	var response model.TCUResponseParameters

	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "f4indexno" {
				err := decoder.DecodeElement(&response.F4IndexNo, &se)
				if err != nil {
					return model.TCUResponseParameters{}, err
				}
			} else if se.Name.Local == "StatusCode" {
				err := decoder.DecodeElement(&response.StatusCode, &se)
				if err != nil {
					return model.TCUResponseParameters{}, err
				}
			} else if se.Name.Local == "StatusDescription" {
				err := decoder.DecodeElement(&response.StatusDescription, &se)
				if err != nil {
					return model.TCUResponseParameters{}, err
				}
			}
		}
	}

	return response, nil
}

func parseXMLMulti(xmlData string) ([]model.TCUResponseParameters, error) {
	var response model.TCUResponse

	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "ResponseParameters" {
				var rp model.TCUResponseParameters
				err := decoder.DecodeElement(&rp, &se)
				if err != nil {
					return nil, err
				}
				response.ResponseParameters = append(response.ResponseParameters, rp)
			}
		}
	}

	return response.ResponseParameters, nil
}

func FindAllNectaSchools() ([]model.NectaSchool, error) {
	db := config.GetDBInstance()
	var schools []model.NectaSchool
	err := db.Model(&model.NectaSchool{}).Find(&schools).Error

	return schools, err
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

func VerifyStudentAccount(apps []model.ApplicantDetail) ([]model.TCUResponseParameters, error) {

	var indexListing = ""

	for _, app := range apps {
		indexListing = indexListing + "<f4indexno>" + app.F4index + "</f4indexno>"
	}

	var request = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
		"<Request>" +
		"<UsernameToken><Username>MNZ</Username>" +
		"<SessionToken>1Hpn63x87qGSRTjr4OfE</SessionToken>" +
		"</UsernameToken><RequestParameters>" + indexListing + "</RequestParameters></Request>"
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

	response, err := parseXMLMulti(string(resp.Body()))

	log.Println(string(resp.Body()))
	log.Printf("response %v", response)
	return response, err

}

func UpdateStudentStatus(db *gorm.DB, students []model.TCUResponseParameters) {

	for _, student := range students {
		db.Where("f4index = ?", student.F4IndexNo).
			Updates(&model.ApplicantDetail{VerificationStatus: student.StatusDescription,
				VerificationCode: strconv.Itoa(student.StatusCode)})
	}

}

func GetStudentResultsBulky(indexNoList []string,examId int)([]model.NectaStudentResult,error ){
	client := resty.New();
	var responResult struct {
		Response []model.NectaStudentResult `json:"response"`
		Status   model.NectaApiResponseStatus `json:"status"`
	}

	var request    struct {
		Particulars []struct {
			IndexNumber string `json:"index_number"`
			ExamYear string `json:"exam_year"`
			ExamId int `json:"exam_id"`
		} `json:"particulars"`
		ApiKey string `json:"api_key"`
	}

	request.ApiKey = "$2y$10$7BFbtDEWB2uac61b96WhlO7tAJp0p4bHbVYxhZgCe.D.WOGgHrG/2"
	
	for _, v := range(indexNoList) {
		v = strings.ToUpper(v)
		if matchIndex(v) {
		indexNo, examYear := splitIndexToParts(v)
		request.Particulars = append(request.Particulars, struct {
			IndexNumber string `json:"index_number"`
			ExamYear string `json:"exam_year"`
			ExamId int `json:"exam_id"`
		}{IndexNumber: indexNo, ExamYear: examYear, ExamId:  examId})
		}else{
			log.Printf("unmatched index number %s",v)
		}
	}

	requestJson, _ := json.Marshal(request)

	resp , err := client.R().
		SetHeader("Content-Type","application/json").
		SetBody(&requestJson).
		SetResult(&responResult).
		Post("https://api.necta.go.tz/api/results/bulk-general")

	if err != nil {
		log.Fatal(err)
	}

	if resp.IsError() {
		log.Fatal(resp.RawResponse)
	}

	return responResult.Response, nil

}

func matchIndex(index string) bool {
    // Define the regular expression pattern
    pattern := `^[SP]\d{4}/\d{4}/\d{4}$`
    
    // Compile the regular expression
    re := regexp.MustCompile(pattern)
    
    // Check if the index matches the pattern
    return re.MatchString(index)
}

func CreateStudentNectaResults(db *gorm.DB, students []model.NectaStudentResult, indexNoList []string, examId int) error {
	// search student results for each index number
	for _, student := range students {
		// update student results
		if student.Status.Code == 1 {
			indexNo := student.Particulars.IndexNumber + "/" + student.Particulars.ExamYear
			indexNo = strings.ReplaceAll(indexNo, "-", "/")
			if (examId == 1){
				cseeResultJson, _ := json.Marshal(student)
				db.Model(&model.ApplicantDetail{}).
				Where("f4index = ?", indexNo).
				Updates(&model.ApplicantDetail{Fname:student.Particulars.FirstName,
					Mname:student.Particulars.MiddleName,
					Lname:student.Particulars.LastName,
					Gender: student.Particulars.Sex,
					F4index: indexNo, CseeResult: string(cseeResultJson)})
			}else{
				acseeResultJson, _ := json.Marshal(student.Results)
				db.Model(&model.ApplicantDetail{}).
				Where("f6index = ?", indexNo).
				Updates(&model.ApplicantDetail{Fname:student.Particulars.FirstName,
					Mname:student.Particulars.MiddleName,
					Lname:student.Particulars.LastName,
					Gender:student.Particulars.Sex,
					AcseeResult: string(acseeResultJson)})
			}
		}else{
			log.Printf("Error updating student %s: %v \n", student.Particulars.IndexNumber, student)
		}
	}
	return nil
}

// split index no exam is S0140/0010/2010 
// split to year and index and replacy foward slash with dash (-)

func splitIndexToParts(indexNo string) (string, string) {
	var index = indexNo[:10]
	var year = indexNo[11:15]
	index = strings.ReplaceAll(index, "/", "-")
	return index, year
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
