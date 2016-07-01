package fetcher

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestFetch(t *testing.T) {
	assert := assert.New(t)
	fetcher := NewTeacherLessonFetcher(nil)
	teacher, _, err := fetcher.Fetch(1523)
	assert.NoError(err)
	assert.Equal("Crizelle", teacher.Name)
}

func TestParseHtml(t *testing.T) {
	assert := assert.New(t)
	fetcher := NewTeacherLessonFetcher(nil)
	file, err := os.Open("testdata/5982.html")
	assert.NoError(err)
	defer file.Close()

	teacher, lessons, err := fetcher.parseHtml(5982, file)
	assert.Equal("Xai", teacher.Name)
	assert.True(len(lessons) > 0)
	//fmt.Printf("%v\n", spew.Sdump(lessons))
}
