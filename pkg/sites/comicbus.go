package sites

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
	"github.com/chromedp/chromedp"
)

type Comicbus struct {
	options *config.Options
}

func NewComicbus(options *config.Options) *Comicbus {
	return &Comicbus{
		options: options,
	}
}

func (c *Comicbus) Initialize(comic *core.Comic) error {
	url := comic.URLSource
	html := getHtmlWithJs(url)
	doc := soup.HTMLParse(html)
	indexes := doc.Find("select", "id", "pageindex").FindAll("option")
	var links []string
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	for i := 1; i <= len(indexes); i++ {
		var html string
		chromedp.Run(ctx,
			chromedp.Navigate(url+"-"+strconv.Itoa(i)),
			chromedp.OuterHTML("html", &html))

		doc := soup.HTMLParse(html)
		imgUrl := "https:" + doc.Find("img", "name", "TheImg").Attrs()["src"]
		fmt.Println(imgUrl)
		links = append(links, imgUrl)
	}
	comic.Links = links
	fmt.Println(links)
	return nil
}

func (c *Comicbus) GetInfo(url string) (string, string) {
	html := getHtmlWithJs(url)
	doc := soup.HTMLParse(html)
	breadcrumb := doc.Find("form").Find("table").Find("td").FullText()
	infos := strings.Split(breadcrumb, ">")
	name := strings.TrimSpace(infos[0])
	issueNumber := strings.TrimSpace(infos[1])
	return name, issueNumber
}

func getHtmlWithJs(url string) string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var html string
	chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &html))

	return html
}

func (c *Comicbus) RetrieveIssueLinks() ([]string, error) {
	url := c.options.Url
	html, ri, err := getHtml(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(html)
	chapters := doc.Find("table", "id", "div_li1").FindAll("a")

	var links []string
	for _, chapter := range chapters {
		suffixUrl, copyright := getAttrs(chapter.Attrs()["onclick"])
		url := getChapterUrl(suffixUrl, ri, copyright)
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Issue Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

func getAttrs(s string) (string, string) {
	s = strings.ReplaceAll(s, "cview(", "")
	s = strings.ReplaceAll(s, ");return false;,", "")
	s = strings.ReplaceAll(s, "'", "")
	attrs := strings.Split(s, ",")
	return strings.TrimSpace(attrs[0]), strings.TrimSpace(attrs[2])
}

func getHtml(url string) (string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var ri string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "RI" {
			ri = cookie.Value
			break
		}
	}

	return string(data), ri, nil
}

func getChapterUrl(url string, ri string, copyright string) string {
	var baseUrl string
	mid := strings.Split(url, "-")[0]
	url = strings.ReplaceAll(url, ".html", "")
	url = strings.ReplaceAll(url, "-", ".html?ch=")

	if ri == "3" && copyright == "1" {
		baseUrl = "https://comic.aya.click/online/c-"
	} else if "17708" == mid {
		baseUrl = "https://comic.aya.click/online/aa-"
	} else {
		baseUrl = "https://comic.aya.click/online/b-"
	}

	return baseUrl + url
}
