package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestComicbusSetup(t *testing.T) {
	opt :=
		&config.Options{
			Url:    "https://www.comicbus.com/html/13313.html",
			All:    true,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	c := NewComicbus(opt)
	comic := new(core.Comic)
	comic.URLSource = "https://comic.aya.click/online/b-11011.html?ch=2"

	err := c.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 24, len(comic.Links))
}

func TestComicbusGetInfo(t *testing.T) {
	opt :=
		&config.Options{
			Url:    "https://www.comicbus.com/html/13313.html",
			All:    true,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	c := NewComicbus(opt)
	name, issueNumber := c.GetInfo("https://comic.aya.click/online/b-13313.html?ch=2")

	assert.Equal(t, "關于我轉生後成為史萊姆的那件事", name)
	assert.Equal(t, "第2集", issueNumber)
}

func TestComicbusRetrieveIssueLinks(t *testing.T) {
	opt :=
		&config.Options{
			Url:    "https://www.comicbus.com/html/13313.html",
			All:    true,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	c := NewComicbus(opt)
	issues, err := c.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 78, len(issues))
}
