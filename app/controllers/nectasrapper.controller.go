package controllers

import (
	"fmt"
	"github.com/OMUHA/oauwebscrapper/app/models"
	"github.com/OMUHA/oauwebscrapper/app/models/necta"
	"github.com/OMUHA/oauwebscrapper/app/repository"
	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func SyncResultToAPI() {

}

func NectaCseeScrapper(ctx *fiber.Ctx) error {
	c := colly.NewCollector(
		colly.AllowedDomains("matokeo.necta.go.tz", "necta.go.tz"),
		colly.Async(true))
	c.UserAgent = "xy"
	c.AllowURLRevisit = true
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   160 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	schoolResultCollector := c.Clone()
	schoolResultCollector.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   120 * time.Second,
			KeepAlive: 60 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	)

	schoolResultCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	schoolResultCollector.OnHTML("body", func(e *colly.HTMLElement) {
		elID := 0
		var students []necta.StudentResult
		var school necta.School
		fmt.Println(e.ChildText("H3"))
		school.Name = strings.Trim(e.ChildText("H3"), "\n")[6:len(strings.Trim(e.ChildText("H3"), "\n"))]
		school.CenterNo = strings.Trim(e.ChildText("H3"), "\n")[:5]
		school.Name = strings.Trim(school.Name, " ")
		school.CenterNo = strings.Trim(school.CenterNo, " ")

		e.ForEach("table", func(_ int, el *colly.HTMLElement) {

			if elID == 2 {
				el.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
					if tr.ChildText("td") == "SEX" {
						return
					}

					var student necta.StudentResult
					tr.ForEach("td", func(id int, td *colly.HTMLElement) {
						switch id {
						case 0:
							student.CandidateNo = strings.Trim(td.Text, "\n")
						case 1:
							student.CandidateGender = strings.Trim(td.Text, "\n")
						case 2:
							student.AggregatePoints = strings.Trim(td.Text, "\n")
						case 3:
							student.ResultDivision = strings.Trim(td.Text, "\n")
						case 4:
							student.ResultsRaw = strings.Trim(td.Text, "\n")
						default:
						}
					})

					student.SchoolName = school.Name
					student.CenterNo = school.CenterNo
					student.ExamType = "csee"
					student.ExamYear = "2022"

					students = append(students, student)
				})
			}

			elID += 1

			if elID == 5 {
				el.ForEach("tr", func(rowID int, tr *colly.HTMLElement) {

					if rowID == 0 {
						tr.ForEach("td", func(colID int, td *colly.HTMLElement) {
							if colID == 1 {
								school.Region = strings.Trim(td.Text, "\n")
							}

						})
					}
				})
			}
		})

		fmt.Println("Importing Data to DB for School : ", e.ChildText("H3"))
		err := repository.StoreSchool(school)
		if err != nil {
			fmt.Println(err)
		}
		err = repository.StoreStudentResults(students)
		if err != nil {
			fmt.Println(err)
		}
	})

	var responseElements []string

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("table", func(_ int, el *colly.HTMLElement) {
			if el.Attr("width") == "700" {
				el.ForEach("a", func(_ int, el *colly.HTMLElement) {
					fmt.Println(el.Attr("href"))
					link := el.Attr("href")
					centerNo := link[8:13]
					if repository.CheckSchoolExists(centerNo) {
						fmt.Println("Center Number already exists ", centerNo)
					}
					err := schoolResultCollector.Visit("https://matokeo.necta.go.tz/csee2022/results/" + centerNo + ".htm")
					if err != nil {
						fmt.Println(err)
						return
					}
					time.Sleep((1 * time.Second) / 4)

				})
			} else {
				return
			}
		})
	})

	for i := 0; i < len(responseElements); i++ {
		fmt.Println(responseElements[i])
	}

	var response models.Response

	response.Data = responseElements
	response.Message = "Success"
	response.Status = 200

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err := c.Visit("https://matokeo.necta.go.tz/csee2022/index.htm")
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response)

}
