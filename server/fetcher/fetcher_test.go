package fetcher

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_           = fmt.Print
	concurrency = flag.Int("concurrency", 1, "concurrency")
	mCountries  *model.MCountries
)

func TestMain(m *testing.M) {
	bootstrap.CheckCLIEnvVars()
	testDBURL := model.ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL())
	var err error
	db, err := model.OpenDB(testDBURL, 1, true) // TODO: env
	if err != nil {
		panic(err)
	}

	mCountries, err = model.NewMCountryService(db).LoadAll()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

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
		Status:     "200 OK",
	}
	resp.Header.Set("Content-Type", "text/html; charset=UTF-8")

	file, err := os.Open("testdata/5982.html")
	if err != nil {
		return nil, err
	}
	resp.Body = file // Close() will be called by client
	return resp, nil
}

type redirectTransport struct{}

func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusFound,
		Status:     "302 Found",
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	resp.Header.Set("Location", "https://twitter.com/")
	return resp, nil
}

func TestFetch(t *testing.T) {
	a := assert.New(t)
	transport := &errorTransport{okThreshold: 0}
	client := &http.Client{Transport: transport}
	fetcher := NewLessonFetcher(client, 1, false, mCountries, nil)
	teacher, _, err := fetcher.Fetch(5982)
	a.Nil(err)
	a.Equal("Xai", teacher.Name)
	a.Equal(1, transport.callCount)
}

//func TestFetchReal(t *testing.T) {
//	a := assert.New(t)
//	http.DefaultClient.Timeout = 10 * time.Second
//	fetcher := NewLessonFetcher(http.DefaultClient, nil)
//	teacher, _, err := fetcher.Fetch(5982)
//	a.Nil(err)
//	a.Equal("Xai", teacher.Name)
//}

func TestFetchRetry(t *testing.T) {
	a := assert.New(t)
	transport := &errorTransport{okThreshold: 2}
	client := &http.Client{Transport: transport}
	fetcher := NewLessonFetcher(client, 1, false, mCountries, nil)
	teacher, _, err := fetcher.Fetch(5982)
	a.Nil(err)
	a.Equal("Xai", teacher.Name)
	a.Equal(2, transport.callCount)
}

func TestFetchRedirect(t *testing.T) {
	a := assert.New(t)
	client := &http.Client{
		Transport:     &redirectTransport{},
		CheckRedirect: redirectErrorFunc,
	}
	fetcher := NewLessonFetcher(client, 1, false, mCountries, nil)
	_, _, err := fetcher.Fetch(5982)
	a.Error(err)
	a.Equal(reflect.TypeOf(&errors.NotFound{}), reflect.TypeOf(err))
}

type responseTransport struct {
	statusCode int
	status     string
	content    string
}

func (t *responseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: t.statusCode,
		Status:     t.status,
		Body:       ioutil.NopCloser(strings.NewReader(t.content)),
	}
	return resp, nil
}

func TestFetchInternalServerError(t *testing.T) {
	a := assert.New(t)
	client := &http.Client{
		Transport: &responseTransport{
			statusCode: http.StatusInternalServerError,
			content:    "Internal Server Error",
		},
	}
	fetcher := NewLessonFetcher(client, 1, false, mCountries, nil)
	_, _, err := fetcher.Fetch(5982)
	a.Error(err)
	a.Contains(err.Error(), "Unknown error in fetchContent")
	a.Contains(err.Error(), "status=500")
}

func TestFetchConcurrency(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	mockTransport, err := NewMockTransport("testdata/5982.html")
	r.NoError(err)
	client := &http.Client{Transport: mockTransport}
	fetcher := NewLessonFetcher(client, *concurrency, false, mCountries, nil)

	const n = 1000
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(teacherID int) {
			defer wg.Done()
			_, _, err := fetcher.Fetch(uint32(teacherID))
			if err != nil {
				fmt.Printf("err = %v\n", err)
				return
			}
		}(i)
	}
	wg.Wait()

	a.Equal(n, mockTransport.NumCalled)
}

func TestParseHTML(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	fetcher := NewLessonFetcher(http.DefaultClient, 1, false, mCountries, nil)
	file, err := os.Open("testdata/5982.html")
	r.NoError(err)
	defer file.Close()

	teacher, lessons, err := fetcher.parseHTML(model.NewTeacher(uint32(5982)), file)
	r.NoError(err)
	a.Equal("Xai", teacher.Name)
	a.Equal(uint16(608), teacher.CountryID) // Philippines
	a.Equal("female", teacher.Gender)
	a.Equal("1980-03-17", teacher.Birthday.Format("2006-01-02"))
	a.Equal(uint32(814), teacher.FavoriteCount)

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
