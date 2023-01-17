package dmm_eikaiwa

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/volatiletech/sqlboiler/v4/types"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/mock"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

var (
	concurrency  = flag.Int("concurrency", 1, "concurrency")
	mCountryList *model2.MCountryList
)

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	testDBURL := model.ReplaceToTestDBURL(nil, config.DefaultVars.DBURL())
	var err error
	db, err := model.OpenDB(testDBURL, 1, config.DefaultVars.DebugSQL)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	mCountries, err := mysql.NewMCountryRepository(db.DB()).FindAll(ctx)
	if err != nil {
		panic(err)
	}
	mCountryList = model2.NewMCountryList(mCountries)

	os.Exit(m.Run())
}

func Test_lessonFetcher_Fetch(t *testing.T) {
	transport := &errorTransport{okThreshold: 0}
	httpClient := &http.Client{Transport: transport}
	fetcher := NewLessonFetcher(httpClient, *concurrency, false, mCountryList, zap.NewNop())
	ctx := context.Background()
	teacher, lessons, err := fetcher.Fetch(ctx, 49393)
	if err != nil {
		t.Fatalf("fetcher.Fetch failed: %v", err)
	}

	assertion.AssertEqual(t, "Judith", teacher.Name, "")
	assertion.AssertEqual(t, 24, len(lessons), "")
	assertion.AssertEqual(t, 1, transport.callCount, "")
}

func Test_lessonFetcher_Fetch_Retry(t *testing.T) {
	transport := &errorTransport{okThreshold: 2}
	client := &http.Client{Transport: transport}
	fetcher := NewLessonFetcher(client, 1, false, mCountryList, zap.NewNop())
	teacher, _, err := fetcher.Fetch(context.Background(), 49393)
	if err != nil {
		t.Fatalf("fetcher.Fetch failed: %v", err)
	}

	assertion.AssertEqual(t, "Judith", teacher.Name, "")
	assertion.AssertEqual(t, 2, transport.callCount, "")
}

func Test_lessonFetcher_Fetch_Redirect(t *testing.T) {
	client := &http.Client{
		Transport:     &redirectTransport{},
		CheckRedirect: redirectErrorFunc,
	}
	fetcher := NewLessonFetcher(client, 1, false, mCountryList, zap.NewNop())
	_, _, err := fetcher.Fetch(context.Background(), 5982)
	if err == nil {
		t.Fatalf("err must not be nil")
	}

	assertion.AssertEqual(t, true, errors.IsNotFound(err), "")
}

func Test_lessonFetcher_Fetch_InternalServerError(t *testing.T) {
	client := &http.Client{
		Transport: &responseTransport{
			statusCode: http.StatusInternalServerError,
			content:    "Internal Server Error",
		},
	}
	fetcher := NewLessonFetcher(client, 1, false, mCountryList, zap.NewNop())
	_, _, err := fetcher.Fetch(context.Background(), 5982)
	if err == nil {
		t.Fatalf("err must not be nil")
	}

	wantTexts := []string{
		"Unknown error in fetchContent",
		"statusCode=500",
	}
	for _, want := range wantTexts {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("err %q doesn't contain text %q", err, want)
		}
	}
}

func Test_lessonFetcher_Fetch_Concurrency(t *testing.T) {
	mockTransport, err := mock.NewHTMLTransport("testdata/49393.html")
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Transport: mockTransport}
	fetcher := NewLessonFetcher(client, *concurrency, false, mCountryList, zap.NewNop())

	const n = 500
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(teacherID int) {
			defer wg.Done()
			_, _, err := fetcher.Fetch(context.Background(), uint(teacherID))
			if err != nil {
				fmt.Printf("err = %v\n", err)
				return
			}
		}(i)
	}
	wg.Wait()

	assertion.AssertEqual(t, n, mockTransport.NumCalled, "")
}

func Test_lessonFetcher_parseHTML(t *testing.T) {
	fetcher := NewLessonFetcher(http.DefaultClient, 1, false, mCountryList, zap.NewNop()).(*lessonFetcher)
	file, err := os.Open("testdata/49393.html")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = file.Close()
	})

	gotTeacher, lessons, err := fetcher.parseHTML(model2.NewTeacher(uint(49393)), file)
	if err != nil {
		t.Fatalf("fetcher.parseHTML failed: %v", err)
	}
	wantTeacher := modeltest.NewTeacher(func(teacher *model2.Teacher) {
		teacher.ID = 49393
		teacher.Name = "Judith"
		teacher.CountryID = int16(608)
		teacher.Birthday = time.Time{}
		teacher.YearsOfExperience = 2
		teacher.FavoriteCount = 403
		teacher.ReviewCount = 0
		teacher.Rating = types.NullDecimal{Big: decimal.New(int64(496), 2)}
		//teacher.LastLessonAt = time.Date(2022, 12, 31, 10, 30, 0, 0, time.UTC)
	})
	assertion.AssertEqual(
		t, wantTeacher, gotTeacher, "",
		cmp.AllowUnexported(decimal.Big{}, big.Int{}),
		cmpopts.IgnoreFields(model2.Teacher{}, "LastLessonAt"),
	)

	assertion.AssertEqual(t, true, len(lessons) > 0, "num of lessons must be greater than zero")
	const dtFormat = "2006-01-02 15:04"
	for _, lesson := range lessons {
		if lesson.Datetime.Format(dtFormat) == "2018-03-01 18:00" {
			assertion.AssertEqual(t, "Finished", lesson.Status, "")
		}
		if lesson.Datetime.Format(dtFormat) == "2018-03-03 06:30" {
			assertion.AssertEqual(t, "Available", lesson.Status, "")
		}
		if lesson.Datetime.Format(dtFormat) == "2018-03-03 02:00" {
			assertion.AssertEqual(t, "Reserved", lesson.Status, "")
		}
	}
}

//<a href="#" class="bt-open" id="a:3:{s:8:&quot;launched&quot;;s:19:&quot;2016-07-01 16:30:00&quot;;s:10:&quot;teacher_id&quot;;s:4:&quot;5982&quot;;s:9:&quot;lesson_id&quot;;s:8:&quot;25880364&quot;;}">予約可</a>

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

	file, err := os.Open("testdata/49393.html")
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
		Body:       io.NopCloser(strings.NewReader("")),
	}
	resp.Header.Set("Location", "https://twitter.com/")
	return resp, nil
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
		Body:       io.NopCloser(strings.NewReader(t.content)),
	}
	return resp, nil
}
