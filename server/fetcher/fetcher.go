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
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/uber-go/zap"
	"gopkg.in/xmlpath.v2"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/601.6.17 (KHTML, like Gecko) Version/9.1.1 Safari/601.6.17"
)

var (
	_              = fmt.Print
	jst            = time.FixedZone("Asia/Tokyo", 9*60*60)
	titleXPath     = xmlpath.MustCompile(`//title`)
	lessonXPath    = xmlpath.MustCompile("//ul[@class='oneday']//li")
	classAttrXPath = xmlpath.MustCompile("@class")
)

type TeacherLessonFetcher struct {
	httpClient *http.Client
	log        zap.Logger
}

func NewTeacherLessonFetcher(httpClient *http.Client, log zap.Logger) *TeacherLessonFetcher {
	if log == nil {
		log = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")))
	}
	return &TeacherLessonFetcher{
		httpClient: httpClient,
		log:        log,
	}
}

func (fetcher *TeacherLessonFetcher) Fetch(teacherId uint32) (*model.Teacher, []*model.Lesson, error) {
	teacher := model.NewTeacher(teacherId)
	var content string
	err := retry.Retry(2, 300*time.Millisecond, func() error {
		var err error
		content, err = fetcher.fetchContent(teacher.Url())
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return fetcher.parseHtml(teacher, strings.NewReader(content))
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

	if resp.StatusCode != 200 {
		return "", errors.Internalf(
			"fetch schedule error: url=%v, status =%v",
			url, resp.StatusCode,
		)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.InternalWrapf(err, "Failed ioutil.ReadAll(): url=%v", url)
	}
	return string(b), nil
}

func (fetcher *TeacherLessonFetcher) parseHtml(
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
		return nil, nil, errors.Internalf("failed to fetch teacher's name: url=%v", teacher.Url)
	}

	dateRegexp := regexp.MustCompile(`([\d]+)月([\d]+)日(.+)`)
	lessons := make([]*model.Lesson, 0, 1000)
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
				originalDate = time.Date(date.Year(), time.Month(month), int(day), 0, 0, 0, 0, jst)
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
			dt := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, jst)
			status := model.LessonStatuses.MustValueForAlias(text)
			fetcher.log.Info(
				"lesson",
				zap.String("dt", dt.Format("2006-01-02 15:04")),
				zap.String("status", model.LessonStatuses.MustName(status)),
			)
			lessons = append(lessons, &model.Lesson{
				TeacherId: teacher.Id,
				Datetime:  dt,
				Status:    model.LessonStatuses.MustName(status),
			})
		} else {
			// nop
		}
	}

	return teacher, lessons, nil
}

func MustInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}
