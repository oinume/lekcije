package fetcher

import (
	"fmt"
	"os"
	"testing"

	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestFetch(t *testing.T) {
	assert := assert.New(t)
	fetcher := NewTeacherLessonFetcher(nil, nil)
	teacher, _, err := fetcher.Fetch(1523)
	assert.NoError(err)
	assert.Equal("Crizelle", teacher.Name)
}

func TestParseHtml(t *testing.T) {
	assert := assert.New(t)
	fetcher := NewTeacherLessonFetcher(nil, nil)
	file, err := os.Open("testdata/5982.html")
	assert.NoError(err)
	defer file.Close()

	teacher, lessons, err := fetcher.parseHtml(model.NewTeacher(uint32(5982)), file)
	assert.Equal("Xai", teacher.Name)
	assert.True(len(lessons) > 0)
	for _, lesson := range lessons {
		if lesson.Datetime.Format("2006-01-02 15:04") == "2016-07-01 11:00" {
			assert.Equal("Finished", lesson.Status)
		}
		if lesson.Datetime.Format("2006-01-02 15:04") == "2016-07-01 16:30" {
			assert.Equal("Reservable", lesson.Status)
		}
		if lesson.Datetime.Format("2006-01-02 15:04") == "2016-07-01 18:00" {
			assert.Equal("Reserved", lesson.Status)
		}
	}
	//fmt.Printf("%v\n", spew.Sdump(lessons))
}

//<a href="#" class="bt-open" id="a:3:{s:8:&quot;launched&quot;;s:19:&quot;2016-07-01 16:30:00&quot;;s:10:&quot;teacher_id&quot;;s:4:&quot;5982&quot;;s:9:&quot;lesson_id&quot;;s:8:&quot;25880364&quot;;}">予約可</a>
