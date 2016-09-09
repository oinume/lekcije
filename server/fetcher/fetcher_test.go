package fetcher

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

type errorTransport struct {
	okThreshold int
	callCount   int
}

func (t *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.callCount++
	if t.callCount < t.okThreshold {
		return nil, fmt.Errorf("Please retry.")
	}

	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	resp.Header.Set("Content-Type", "text/html; charset=UTF-8")

	file, err := os.Open("testdata/5982.html")
	if err != nil {
		return nil, err
	}
	resp.Body = file // Close() will be called by client
	return resp, nil
}

func TestFetch(t *testing.T) {
	a := assert.New(t)
	transport := &errorTransport{okThreshold: 0}
	client := &http.Client{Transport: transport}
	fetcher := NewTeacherLessonFetcher(client, nil)
	teacher, _, err := fetcher.Fetch(5982)
	a.NoError(err)
	a.Equal("Xai", teacher.Name)
	a.Equal(1, transport.callCount)
}

//func TestFetchReal(t *testing.T) {
//	a := assert.New(t)
//	http.DefaultClient.Timeout = 10 * time.Second
//	fetcher := NewTeacherLessonFetcher(http.DefaultClient, nil)
//	teacher, _, err := fetcher.Fetch(5982)
//	a.NoError(err)
//	a.Equal("Xai", teacher.Name)
//}

func TestFetchRetry(t *testing.T) {
	a := assert.New(t)
	transport := &errorTransport{okThreshold: 2}
	client := &http.Client{Transport: transport}
	fetcher := NewTeacherLessonFetcher(client, nil)
	teacher, _, err := fetcher.Fetch(5982)
	a.NoError(err)
	a.Equal("Xai", teacher.Name)
	a.Equal(2, transport.callCount)
}

func TestParseHtml(t *testing.T) {
	a := assert.New(t)
	fetcher := NewTeacherLessonFetcher(http.DefaultClient, nil)
	file, err := os.Open("testdata/5982.html")
	a.NoError(err)
	defer file.Close()

	teacher, lessons, err := fetcher.parseHtml(model.NewTeacher(uint32(5982)), file)
	a.Equal("Xai", teacher.Name)
	a.True(len(lessons) > 0)
	for _, lesson := range lessons {
		if lesson.Datetime.Format("2006-01-02 15:04") == "2016-07-01 11:00" {
			a.Equal("Finished", lesson.Status)
		}
		if lesson.Datetime.Format("2006-01-02 15:04") == "2016-07-01 16:30" {
			a.Equal("Available", lesson.Status)
		}
		if lesson.Datetime.Format("2006-01-02 15:04") == "2016-07-01 18:00" {
			a.Equal("Reserved", lesson.Status)
		}
	}
	//fmt.Printf("%v\n", spew.Sdump(lessons))
}

//<a href="#" class="bt-open" id="a:3:{s:8:&quot;launched&quot;;s:19:&quot;2016-07-01 16:30:00&quot;;s:10:&quot;teacher_id&quot;;s:4:&quot;5982&quot;;s:9:&quot;lesson_id&quot;;s:8:&quot;25880364&quot;;}">予約可</a>
