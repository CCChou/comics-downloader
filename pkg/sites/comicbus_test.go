package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCartoonmadRetrieveIssueLinks(t *testing.T) {
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
