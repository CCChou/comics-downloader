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
	options       *config.Options
	chromeContext context.Context
	chromeCancel  context.CancelFunc
}

func NewComicbus(options *config.Options) *Comicbus {
	return &Comicbus{
		options: options,
	}
}

func (c *Comicbus) Initialize(comic *core.Comic) error {
	url := comic.URLSource
	html := c.getHtmlWithJs(url)
	doc := soup.HTMLParse(html)
	indexes := doc.Find("select", "id", "pageindex").FindAll("option")
	var links []string
	for i := 1; i <= len(indexes); i++ {
		html := c.getHtmlWithJs(url + "-" + strconv.Itoa(i))
		doc := soup.HTMLParse(html)
		imgUrl := "https:" + doc.Find("img", "name", "TheImg").Attrs()["src"]
		links = append(links, imgUrl)
	}
	comic.Links = links

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Links found: %s", strings.Join(links, " ")))
	}

	return nil
}

func (c *Comicbus) GetInfo(url string) (string, string) {
	html := c.getHtmlWithJs(url)
	doc := soup.HTMLParse(html)
	breadcrumb := doc.Find("form").Find("table").Find("td").FullText()
	infos := strings.Split(breadcrumb, ">")
	name := strings.TrimSpace(infos[0])
	issueNumber := strings.TrimSpace(infos[1])

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Name: %s, Issue Number: %s", name, issueNumber))
	}

	return name, issueNumber
}

func (c *Comicbus) getHtmlWithJs(url string) string {
	ctx := c.getContext()
	var html string
	chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &html))

	return html
}

func (c *Comicbus) getContext() context.Context {
	if c.chromeContext == nil {
		ctx, cancel := chromedp.NewContext(context.Background())
		c.chromeContext = ctx
		c.chromeCancel = cancel
	}
	return c.chromeContext
}

func (c *Comicbus) RetrieveIssueLinks() ([]string, error) {
	url := c.options.Url
	html, ri, err := c.getHtml(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(html)
	chapters := doc.Find("table", "id", "div_li1").FindAll("a")

	var links []string
	for _, chapter := range chapters {
		suffixUrl, copyright := c.getAttrs(chapter.Attrs()["onclick"])
		url := c.getChapterUrl(suffixUrl, ri, copyright)
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Issue Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

func (c *Comicbus) getAttrs(s string) (string, string) {
	s = strings.ReplaceAll(s, "cview(", "")
	s = strings.ReplaceAll(s, ");return false;,", "")
	s = strings.ReplaceAll(s, "'", "")
	attrs := strings.Split(s, ",")
	return strings.TrimSpace(attrs[0]), strings.TrimSpace(attrs[2])
}

func (c *Comicbus) getHtml(url string) (string, string, error) {
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

func (c *Comicbus) getChapterUrl(url string, ri string, copyright string) string {
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
