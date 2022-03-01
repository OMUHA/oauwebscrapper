 package necta

 import "github.com/gocolly/colly"

 type CseefrontPage struct {
	PageID   int    `json:"pageID"`
	PageName string `json:"pageName"`
	ElementDetail colly.HTMLElement
}

type ListedUrl struct {
	Url string `json:"url"`
}
