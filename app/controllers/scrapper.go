package controllers

import (
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

var endCount = 100 //255731
func DownloadAppData(ctx *fiber.Ctx) error {
	c := colly.NewCollector(
		colly.AllowedDomains("uims.tcu.go.tz", "tcu.go.tz"),
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
	})

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	c.SetCookieJar(jar)

	cookie := &http.Cookie{
		Name:   "PHPSESSID",
		Value:  "74oi0du5vvdvd5o62b5hjmuch5",
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

					time.Sleep((1 * time.Second) / 4)

				})
			} else {
				return
			}
		})
	})
	var response models.Response

	response.Data = nil
	response.Message = "Success"
	response.Status = 200

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	startCount := 1

	repeated := endCount / 20000
	repeated += 1
	startAt := startCount
	for ix := startCount; ix <= repeated; ix++ {
		fmt.Printf("Downloading from: %d to %d\n", startAt, startAt+20000)
		go anotherGoFuncToDownload(c.Clone(), startAt, startAt+20000)
		startAt = ix * 20000

	}

	return ctx.Status(200).JSON(response)

}

func anotherGoFuncToDownload(schoolResultCollector *colly.Collector, start, end int) {
	if start > endCount {
		fmt.Println("finishied downloading")
		panic("program must end")
	} else {
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

			go repository.CreateStudentDetails(db, student)
		})

		for i := start; i <= end; i++ {
			err := schoolResultCollector.Visit("https://uims.tcu.go.tz/index.php?r=selectedApplicantsUploadedThroughApi/view&id=" + strconv.Itoa(i))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
