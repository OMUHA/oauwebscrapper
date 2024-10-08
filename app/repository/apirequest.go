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
	
	var responResult struct {
		Candidates []model.NectaStudentResult `json:"candidates"`
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

	//requestJson, _ := json.Marshal(request)

	client := resty.New();
	// Add a request hook to log requests
    /* client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
        log.Printf("Request URL: %s", r.URL)
        log.Printf("Request Method: %s", r.Method)
        log.Printf("Request Headers: %v", r.Header)
        log.Printf("Request Body: %s", r.Body)
        return nil
    }) */

    // Add a response hook to log responses
     /* client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
        log.Printf("Response Status Code: %d", r.StatusCode())
        log.Printf("Response Headers: %v", r.Header())
        log.Printf("Response Body: %s", r.String())
        return nil
    }) */
	resp , err := client.R().
		SetHeader("Content-Type","application/json").
		SetBody(request).
		SetResult(&responResult).
		Post("https://api.necta.go.tz/api/results/bulk-general")

	if err != nil {
		log.Fatal(err)
	}

	if resp.IsError() {
		log.Fatal(resp.RawResponse)
	}

	log.Printf("Response %+v", responResult.Status)
	return responResult.Candidates, nil

}

func matchIndex(index string) bool {
    // Define the regular expression pattern
    pattern := `^[SP]\d{4}/\d{4}/\d{4}$`
    
    // Compile the regular expression
    re := regexp.MustCompile(pattern)
    
    // Check if the index matches the pattern
    return re.MatchString(index)
}

func mapIndexFromList(index string, indexList []string) string {
	index = strings.ReplaceAll(index, "-", "/")
	for _, v := range indexList {
		if strings.Contains( v,index) {
			return v
		}
	}
	return index
}

func mergeSubjectsToString(subjects []model.Subject) string {
	var mappedString = ""
	for _, subject := range subjects {
		mappedString = mappedString + subject.SubjectName + " - " + subject.Grade + ","
	}
	return mappedString
}

func CreateStudentNectaResults(db *gorm.DB, students []model.NectaStudentResult, indexNoList []string, examId int) error {
	// search student results for each index number
	log.Printf("Total students %d \n", len(students))
	for _, student := range students {
		// update student results
		if student.Status.Code == 1 {
			indexNo := mapIndexFromList(student.Particulars.IndexNumber,indexNoList)
			if (examId == 1){
				cseeResultJson, _ := json.Marshal(student.Subjects)
				mappedString := mergeSubjectsToString(student.Subjects)
				err := db.Model(&model.ApplicantDetail{}).
				Where("f4index = ?", indexNo).
				Updates(&model.ApplicantDetail{Fname:student.Particulars.FirstName,
					Mname:student.Particulars.MiddleName,
					Lname:student.Particulars.LastName,
					Gender: student.Particulars.Sex,
					F4result: mappedString,
					CseeCenterName: student.Particulars.CenterName,
					CseeDivision: student.Results.Division,
					CseePoints: student.Results.Points,
					F4index: indexNo, CseeResult: string(cseeResultJson)}).Error 
				if err != nil {
					log.Printf("Error updating csee student %s:   %s\n", indexNo,err.Error())
				}
			}else{
				acseeResultJson, _ := json.Marshal(student.Subjects)
				mappedString := mergeSubjectsToString(student.Subjects)
				err := db.Model(&model.ApplicantDetail{}).
				Where("f6_index = ?", indexNo).
				Updates(&model.ApplicantDetail{Fname:student.Particulars.FirstName,
					Mname:student.Particulars.MiddleName,
					Lname:student.Particulars.LastName,
					Gender:student.Particulars.Sex,
					F6Result: mappedString,
					AcseeCenterName:student.Particulars.CenterName,
					AcseeDivision:student.Results.Division,
					AcseePoints:student.Results.Points,
					AcseeResult: string(acseeResultJson)}).Error
					if err != nil {
						log.Printf("Error updating acsee student %s:   %s\n", indexNo,err.Error())
					}
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
