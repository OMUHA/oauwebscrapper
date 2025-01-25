package controllers

import (
	"fmt"
	"github.com/OMUHA/oauwebscrapper/app/models"
	"github.com/OMUHA/oauwebscrapper/app/models/necta"
	"github.com/OMUHA/oauwebscrapper/app/repository"
	"github.com/OMUHA/oauwebscrapper/config"
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
		colly.AllowedDomains("onlinesys.necta.go.tz", "matokeo.necta.go.tz", "necta.go.tz"),
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

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	yearID := ctx.Params("yearID")

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
		school.CenterNo = strings.ToUpper(strings.Trim(school.CenterNo, " "))

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
					student.ExamYear = yearID
					student.IndexNo = student.CandidateNo + "/" + yearID

					grades := repository.ExtractGrades(student.ResultsRaw)
					for subject, grade := range grades {
						switch subject {
						case "B/MATH":
							student.Bmath = grade
							student.BmathPts = repository.CalculateGradePoint(grade)
							break
						case "ENGL":
							student.Eng = grade
							student.EngPts = repository.CalculateGradePoint(grade)
							break
						case "BIO":
							student.Bio = grade
							student.BioPts = repository.CalculateGradePoint(grade)
							break
						case "PHY":
							student.Phy = grade
							student.PhyPts = repository.CalculateGradePoint(grade)
							break
						case "CHEM":
							student.Chem = grade
							student.ChemPts = repository.CalculateGradePoint(grade)
						}
					}
					if strings.Contains(strings.ToUpper(student.CandidateNo), "S") {
						student.CandidateType = "S"
					} else {
						student.CandidateType = "P"
					}

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
		go func() {
			err := repository.StoreSchool(school)
			if err != nil {
				fmt.Println(err)
			}
			err = repository.StoreStudentResults(students)
			if err != nil {
				fmt.Println(err)
			}
		}()
	})

	var responseElements []string

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("table", func(_ int, el *colly.HTMLElement) {
			if el.Attr("width") == "700" {
				el.ForEach("a", func(_ int, el *colly.HTMLElement) {
					fmt.Println(el.Attr("href"))
					link := el.Attr("href")
					centerNo := ""
					if yearID == "2023" {
						centerNo = link[17:22]
					} else {
						centerNo = link[8:13]
					}

					log.Println(centerNo)
					log.Println(link)
					centerNoOg := centerNo
					centerNo = strings.ToUpper(centerNo)
					if repository.CheckSchoolExists(centerNo) && repository.CheckNectaSchoolHasStudents(centerNo, yearID, "csee") {
						fmt.Println("Center Number already exists ", centerNo)
						return
					} else {
						if yearID == "2024" {
							err := schoolResultCollector.Visit("https://matokeo.necta.go.tz/results/2024/csee/CSEE2024/CSEE2024/results/" + centerNoOg + ".htm")
							if err != nil {
								fmt.Println(err)
								return
							}
						} else if yearID == "2023" {
							err := schoolResultCollector.Visit("https://matokeo.necta.go.tz/results/2023/csee/CSEE2023/results/" + centerNoOg + ".htm")
							if err != nil {
								fmt.Println(err)
								return
							}
						} else if yearID == "2022" {
							err := schoolResultCollector.Visit("https://onlinesys.necta.go.tz/results/2022/csee/results/" + centerNoOg + ".htm")
							if err != nil {
								fmt.Println(err)
								return
							}
						} else {
							err := schoolResultCollector.Visit("https://onlinesys.necta.go.tz/results/" + yearID + "/csee/results/" + centerNoOg + ".htm")
							if err != nil {
								fmt.Println(err)
								return
							}
						}
						//time.Sleep(1 * time.Second)
					}
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

	if yearID == "2024" {
		err := c.Visit("https://matokeo.necta.go.tz/results/2024/csee/CSEE2024/CSEE2024.htm")
		if err != nil {
			return err
		}
	} else if yearID == "2023" {
		err := c.Visit("https://matokeo.necta.go.tz/results/2023/csee/CSEE%202023.htm")
		if err != nil {
			return err
		}
	} else if yearID == "2022" {
		err := c.Visit("https://onlinesys.necta.go.tz/results/2022/csee/index.htm")
		if err != nil {
			return err
		}
	} else {
		err := c.Visit("https://onlinesys.necta.go.tz/results/" + yearID + "/csee/csee.htm")
		if err != nil {
			return err
		}
	}

	return ctx.Status(200).JSON(response)

}

func NectaACseeScrapper(ctx *fiber.Ctx) error {
	c := colly.NewCollector(
		colly.AllowedDomains("onlinesys.necta.go.tz", "matokeo.necta.go.tz", "necta.go.tz"),
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

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	yearID := ctx.Params("yearID")

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
		school.CenterNo = strings.ToUpper(strings.Trim(school.CenterNo, " "))

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
					student.ExamType = "acsee"
					student.ExamYear = yearID
					student.IndexNo = student.CandidateNo + "/" + yearID
					if strings.Contains(strings.ToUpper(student.CandidateNo), "S") {
						student.CandidateType = "S"
					} else {
						student.CandidateType = "P"
					}
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
		go func() {
			err := repository.StoreSchool(school)
			if err != nil {
				fmt.Println(err)
			}
			err = repository.StoreStudentResults(students)
			if err != nil {
				fmt.Println(err)
			}
		}()
	})

	var responseElements []string

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("table", func(_ int, el *colly.HTMLElement) {
			if el.Attr("width") == "700" {
				el.ForEach("a", func(_ int, el *colly.HTMLElement) {
					fmt.Println(el.Attr("href"))
					link := el.Attr("href")
					centerNo := ""
					centerNo = link[8:13]
					log.Println(centerNo)
					log.Println(link)
					centerNoOg := centerNo
					centerNo = strings.ToUpper(centerNo)
					if repository.CheckSchoolExists(centerNo) && repository.CheckNectaSchoolHasStudents(centerNo, yearID, "acsee") {
						fmt.Println("Center Number already exists ", centerNo)
						return
					} else {
						if yearID == "2023" {
							err := schoolResultCollector.Visit("https://matokeo.necta.go.tz/results/2023/acsee/results/" + centerNoOg + ".htm")
							if err != nil {
								fmt.Println(err)
								return
							}
						} else {
							err := schoolResultCollector.Visit("https://onlinesys.necta.go.tz/results/" + yearID + "/acsee/results/" + centerNoOg + ".htm")
							if err != nil {
								fmt.Println(err)
								return
							}
						}
						//time.Sleep(1 * time.Second)
					}
				})
			} else {
				log.Println("FAiled to get body content")
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

	if yearID == "2023" {
		err := c.Visit("https://matokeo.necta.go.tz/results/2023/acsee/index.htm")
		if err != nil {
			return err
		}
	} else {
		err := c.Visit("https://onlinesys.necta.go.tz/results/" + yearID + "/acsee/index.htm")
		if err != nil {
			return err
		}
	}

	return ctx.Status(200).JSON(response)

}

func NectaUpdateStudent(ctx *fiber.Ctx) error {

	var response models.Response
	response.Data = nil
	response.Message = "Process started"
	processStudents()
	return ctx.Status(200).JSON(response)
}

func processStudents() {
	db := config.GetDBInstance()
	var yearsList []int
	/*yearsList = append(yearsList, 2020)*/
	yearsList = append(yearsList, 2021)
	yearsList = append(yearsList, 2022)
	yearsList = append(yearsList, 2023)
	for _, year := range yearsList {
		// Retrieve batch of records from database
		repository.UpdateStudentResults(db, year)

	}

	log.Println("All records updated successfully")
}
