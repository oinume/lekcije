package fetcher

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/retry"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/uber-go/zap"
	"gopkg.in/xmlpath.v2"
)

const (
	userAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html"
)

var redirectErrorFunc = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

var (
	_                 = fmt.Print
	fetcherHTTPClient = &http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: redirectErrorFunc,
	}
	titleXPath      = xmlpath.MustCompile(`//title`)
	attributesXPath = xmlpath.MustCompile(`//div[@class='confirm low']/dl`)
	lessonXPath     = xmlpath.MustCompile(`//ul[@class='oneday']//li`)
	classAttrXPath  = xmlpath.MustCompile(`@class`)
)

type TeacherLessonFetcher struct {
	httpClient *http.Client
	semaphore  chan struct{}
	log        zap.Logger
}

func NewTeacherLessonFetcher(httpClient *http.Client, concurrency int, log zap.Logger) *TeacherLessonFetcher {
	if httpClient == nil {
		httpClient = fetcherHTTPClient
	}
	if concurrency < 1 {
		concurrency = 1
	}
	if log == nil {
		log = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")))
	}
	semaphore := make(chan struct{}, concurrency)
	return &TeacherLessonFetcher{
		httpClient: httpClient,
		semaphore:  semaphore,
		log:        log,
	}
}

func (fetcher *TeacherLessonFetcher) Fetch(teacherID uint32) (*model.Teacher, []*model.Lesson, error) {
	fetcher.semaphore <- struct{}{}
	defer func() {
		<-fetcher.semaphore
	}()

	teacher := model.NewTeacher(teacherID)
	var content string
	err := retry.Retry(2, 300*time.Millisecond, func() error {
		var err error
		content, err = fetcher.fetchContent(teacher.URL())
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return fetcher.parseHTML(teacher, strings.NewReader(content))
}

func (fetcher *TeacherLessonFetcher) fetchContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.InternalWrapf(err, "Failed to create HTTP request: url=%v", url)
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := fetcher.httpClient.Do(req)
	if err != nil {
		return "", errors.InternalWrapf(err, "Failed httpClient.Do(): url=%v", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound {
		return "", errors.NotFoundf("Teacher not found: url=%v, status=%v", url, resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Internalf(
			"fetchContent error: url=%v, status=%v",
			url, resp.StatusCode,
		)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.InternalWrapf(err, "Failed ioutil.ReadAll(): url=%v", url)
	}
	return string(b), nil
}

func (fetcher *TeacherLessonFetcher) parseHTML(
	teacher *model.Teacher,
	html io.Reader,
) (*model.Teacher, []*model.Lesson, error) {
	root, err := xmlpath.ParseHTML(html)
	if err != nil {
		return nil, nil, err
	}

	// teacher name
	if title, ok := titleXPath.String(root); ok {
		teacher.Name = strings.Trim(strings.Split(title, "-")[0], " ")
	} else {
		return nil, nil, errors.Internalf("failed to fetch teacher's name: url=%v", teacher.URL)
	}

	nameXPath := xmlpath.MustCompile(`./dt`)
	valueXPath := xmlpath.MustCompile(`./dd`)
	for iter := attributesXPath.Iter(root); iter.Next(); {
		node := iter.Node()
		if name, ok := nameXPath.String(node); ok {
			if value, ok := valueXPath.String(node); ok {
				fetcher.setAttribute(teacher, strings.TrimSpace(name), strings.TrimSpace(value))
				//fmt.Printf("name = %v, value = %v\n", strings.TrimSpace(name), strings.TrimSpace(value))
			}
		}
	}
	//fmt.Printf("teacher = %+v\n", teacher)

	dateRegexp := regexp.MustCompile(`([\d]+)月([\d]+)日(.+)`)
	lessons := make([]*model.Lesson, 0, 1000)
	now := time.Now()
	originalDate := time.Now().Truncate(24 * time.Hour)
	date := originalDate
	// lessons
	for iter := lessonXPath.Iter(root); iter.Next(); {
		node := iter.Node()
		timeClass, ok := classAttrXPath.String(node)
		if !ok {
			continue
		}

		text := strings.Trim(node.String(), " ")

		//fmt.Printf("text = '%v', timeClass = '%v'\n", text, timeClass)
		fetcher.log.Debug("Scraping as", zap.String("timeClass", timeClass), zap.String("text", text))

		// blank, available, reserved
		if timeClass == "date" {
			group := dateRegexp.FindStringSubmatch(text)
			if len(group) > 0 {
				month, day := MustInt(group[1]), MustInt(group[2])
				year := date.Year()
				if now.Month() == time.December && month == 1 {
					year = now.Year() + 1
				}
				originalDate = time.Date(
					year, time.Month(month), int(day),
					0, 0, 0, 0,
					config.LocalTimezone(),
				)
				date = originalDate
			}
		} else if strings.HasPrefix(timeClass, "t-") && text != "" {
			tmp := strings.Split(timeClass, "-")
			hour, minute := MustInt(tmp[1]), MustInt(tmp[2])
			if hour >= 24 {
				// Convert 24:30 -> 00:30
				hour -= 24
				if date.Unix() == originalDate.Unix() {
					// Set date to next day for 24:30
					date = date.Add(24 * time.Hour)
				}
			}
			dt := time.Date(
				date.Year(), date.Month(), date.Day(),
				hour, minute, 0, 0,
				config.LocalTimezone(),
			)
			status := model.LessonStatuses.MustValueForAlias(text)
			fetcher.log.Debug(
				"lesson",
				zap.String("dt", dt.Format("2006-01-02 15:04")),
				zap.String("status", model.LessonStatuses.MustName(status)),
			)
			lessons = append(lessons, &model.Lesson{
				TeacherID: teacher.ID,
				Datetime:  dt,
				Status:    model.LessonStatuses.MustName(status),
			})
		} else {
			// nop
		}
	}

	return teacher, lessons, nil
}

func (fetcher *TeacherLessonFetcher) setAttribute(teacher *model.Teacher, name string, value string) error {
	switch name {
	case "国籍":
		teacher.Nationality = value
	case "誕生日":
		t, err := time.Parse("2006-01-02", value)
		if err != nil {
			return err
		}
		teacher.Birthday = t
	case "性別":
		teacher.Gender = value
	case "経歴":
		teacher.YearsOfExperience = value
	}
	return nil
}

func (fetcher *TeacherLessonFetcher) Close() {
	close(fetcher.semaphore)
}

func MustInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}
