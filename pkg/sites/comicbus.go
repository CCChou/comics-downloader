package sites

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

type Comicbus struct {
	options *config.Options
}

func NewComicbus(options *config.Options) *Comicbus {
	return &Comicbus{
		options: options,
	}
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

	fmt.Println(links)

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
