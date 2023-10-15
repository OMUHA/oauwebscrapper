package controllers

import (
	"crypto/tls"
	"fmt"
	"github.com/OMUHA/oauwebscrapper/app/model"
	"github.com/OMUHA/oauwebscrapper/app/models"
	"github.com/OMUHA/oauwebscrapper/app/repository"
	"github.com/OMUHA/oauwebscrapper/config"
	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var endCount = 464211
var studentsLimit = 10000
var startCount = 264211

func VerifyStudentList(ctx *fiber.Ctx) error {

	db := config.GetDBInstance()
	limitStudent := 10000
	totalEntries := int(repository.GetTotalStuentDetaisl(db))
	totalGroups := (totalEntries / 10000) + 1

	var startFilter = 0

	log.Printf("Verifying %d", totalEntries)
	for i := 0; i < totalGroups; i++ {
		var students = repository.GetApplicantDataLimited(db, startFilter, limitStudent)
		log.Printf("Student %d ", len(students))

		if len(students) > 0 {
			status, err := repository.VerifyStudentAccount(students)

			if err != nil {
				log.Printf("student  error %s", err.Error())
			}

			repository.UpdateStudentStatus(db, status)
			startFilter = startFilter + limitStudent
		}
	}

	var response models.Response
	response.Data = nil
	response.Message = "Success"
	response.Status = 200
	return ctx.Status(200).JSON(response)
}

func DownloadAppData(ctx *fiber.Ctx) error {
	c := colly.NewCollector(
		colly.AllowedDomains("uims.tcu.go.tz", "tcu.go.tz"),
		colly.Async(true))
	c.UserAgent = "deio"
	c.AllowURLRevisit = true
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // This will disable SSL verification
		},
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   160 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 30 * time.Second,
	})

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	c.SetCookieJar(jar)

	cookie := &http.Cookie{
		Name:   "PHPSESSID",
		Value:  "ob20vfsdr7e90r2615a7ts08e0",
		Domain: "uims.tcu.go.tz",
	}

	u := "https://uims.tcu.go.tz/"
	c.Cookies(u)
	l, _ := url.Parse(u)
	var cookies []*http.Cookie
	cookies = append(cookies, cookie)

	jar.SetCookies(l, cookies)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	var response models.Response

	response.Data = nil
	response.Message = "Success"
	response.Status = 200

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	repeated := ((endCount - startCount) + 1) / studentsLimit
	repeated += 1
	startAt := startCount
	log.Printf("Starting Downloading %d", startAt)

	for ix := 1; ix <= repeated; ix++ {
		fmt.Printf("Downloading from: %d to %d\n", startAt, startAt+studentsLimit)
		go anotherGoFuncToDownload(c.Clone(), startAt, startAt+studentsLimit)
		startAt = (ix * studentsLimit) + startCount
	}

	return ctx.Status(200).JSON(response)

}

func anotherGoFuncToDownload(schoolResultCollector *colly.Collector, start, end int) {
	if start > endCount {
		fmt.Println("finishied downloading")
		panic("program must end")
	} else {

		schoolResultCollector.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		db := config.GetDBInstance()

		schoolResultCollector.OnHTML("table.detail-view tbody", func(e *colly.HTMLElement) {
			var student model.ApplicantDetail

			e.ForEach("tr", func(id int, el *colly.HTMLElement) {
				switch id {
				case 0:
					student.HliID = strings.Trim(el.ChildText("td"), " ")
					break
				case 1:
					student.F4index = strings.Trim(el.ChildText("td"), " ")
					break
				case 2:
					student.F6Index = strings.Trim(el.ChildText("td"), " ")
					break
				case 3:
					student.Programs = strings.Trim(el.ChildText("td"), " ")
					break
				case 4:
					student.MobileNumber = strings.Trim(el.ChildText("td"), " ")
					break
				case 5:
					student.EmailAddress = strings.Trim(el.ChildText("td"), " ")
					break
				case 6:
					student.AdmissionStatus = strings.Trim(el.ChildText("td"), " ")
					break
				case 7:
					student.AdmittedProgram = strings.Trim(el.ChildText("td"), " ")
					break
				case 8:
					student.Comment = strings.Trim(el.ChildText("td"), " ")
					break
				}
			})
			fmt.Printf("extracted %s %s %s \n", student.HliID, student.F4index, student.F6Index)
			go repository.CreateStudentDetails(db, student)
		})

		schoolResultCollector.OnError(func(r *colly.Response, err error) {
			fmt.Printf("error %s", err)
			fmt.Printf("TCUResponse %s ", r.Body)
		})

		for i := start; i <= end; i++ {
			err := schoolResultCollector.Visit("https://uims.tcu.go.tz/index.php?r=selectedApplicantsUploadedThroughApi/view&id=" + strconv.Itoa(i))
			if err != nil {
				fmt.Println(err.Error())
			}
			time.Sleep((2 * time.Second) / 5)
		}
	}
}

func filterStudentsf(ctx *fiber.Ctx) error {
	var notselectedFiler []string

	notselectedFiler = append(notselectedFiler, "provisionaladmission", "Admitted", "")

	var response models.Response

	return ctx.Status(200).JSON(response)
}
