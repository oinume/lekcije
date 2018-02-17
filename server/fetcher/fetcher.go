package fetcher

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"os"

	"github.com/Songmu/retry"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/text/width"
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
	defaultHTTPClient = &http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: redirectErrorFunc,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10, // TODO: use `concurrency`
			Proxy:               http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	titleXPath      = xmlpath.MustCompile(`//title`)
	attributesXPath = xmlpath.MustCompile(`//div[@class='confirm low']/dl`)
	lessonXPath     = xmlpath.MustCompile(`//ul[@class='oneday']//li`)
	classAttrXPath  = xmlpath.MustCompile(`@class`)
)

type teacherLessons struct {
	teacher *model.Teacher
	lessons []*model.Lesson
}

type LessonFetcher struct {
	httpClient *http.Client
	semaphore  chan struct{}
	caching    bool
	cache      map[uint32]*teacherLessons
	cacheLock  *sync.RWMutex
	logger     *zap.Logger
	mCountries *model.MCountries
}

func NewLessonFetcher(
	httpClient *http.Client, concurrency int, caching bool,
	mCountries *model.MCountries, log *zap.Logger,
) *LessonFetcher {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	if concurrency < 1 {
		concurrency = 1
	}
	if log == nil {
		log = logger.NewZapLogger(nil, []io.Writer{os.Stderr}, zapcore.InfoLevel)
	}
	semaphore := make(chan struct{}, concurrency)
	cache := make(map[uint32]*teacherLessons, 5000)
	return &LessonFetcher{
		httpClient: httpClient,
		semaphore:  semaphore,
		caching:    caching,
		cache:      cache,
		cacheLock:  &sync.RWMutex{},
		logger:     log,
		mCountries: mCountries,
	}
}

func (fetcher *LessonFetcher) Fetch(teacherID uint32) (*model.Teacher, []*model.Lesson, error) {
	fetcher.semaphore <- struct{}{}
	defer func() {
		<-fetcher.semaphore
	}()

	// Check cache
	if fetcher.caching {
		fetcher.cacheLock.RLock()
		if c, ok := fetcher.cache[teacherID]; ok {
			fetcher.cacheLock.RUnlock()
			return c.teacher, c.lessons, nil
		}
		fetcher.cacheLock.RUnlock()
	}

	teacher := model.NewTeacher(teacherID)
	var content io.ReadCloser
	err := retry.Retry(2, 300*time.Millisecond, func() error {
		var err error
		content, err = fetcher.fetchContent(teacher.URL())
		return err
	})
	defer content.Close()
	if err != nil {
		return nil, nil, err
	}
	return fetcher.parseHTML(teacher, content)
}

func (fetcher *LessonFetcher) fetchContent(url string) (io.ReadCloser, error) {
	nopCloser := ioutil.NopCloser(strings.NewReader(""))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nopCloser, errors.InternalWrapf(err, "Failed to create HTTP request: url=%v", url)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := fetcher.httpClient.Do(req)
	if err != nil {
		return nopCloser, errors.InternalWrapf(err, "Failed httpClient.Do(): url=%v", url)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.Body, nil
	case http.StatusMovedPermanently, http.StatusFound:
		_ = resp.Body.Close()
		return nopCloser, errors.NotFoundf("Teacher not found: url=%v, status=%v", url, resp.StatusCode)
	default:
		body, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nopCloser, errors.Internalf(
			"Unknown error in fetchContent: url=%v, status=%v, body=%v",
			url, resp.StatusCode, string(body),
		)
	}
}

func (fetcher *LessonFetcher) parseHTML(
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

	// Nationality, birthday, etc...
	fetcher.parseTeacherAttribute(teacher, root)
	// FavoriteCount
	fetcher.parseTeacherFavoriteCount(teacher, root)

	dateRegexp := regexp.MustCompile(`([\d]+)月([\d]+)日(.+)`)
	lessons := make([]*model.Lesson, 0, 1000)
	now := time.Now().In(config.LocalTimezone())
	originalDate := time.Now().In(config.LocalTimezone()).Truncate(24 * time.Hour)
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
		fetcher.logger.Debug("Scraping as", zap.String("timeClass", timeClass), zap.String("text", text))

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
			fetcher.logger.Debug(
				"lesson",
				zap.String("dt", dt.Format("2006-01-02 15:04")),
				zap.String("status", model.LessonStatuses.MustName(status)),
			)
			lessons = append(lessons, &model.Lesson{
				TeacherID: teacher.ID,
				Datetime:  dt,
				Status:    model.LessonStatuses.MustName(status),
			})
		}
		// TODO: else
	}

	// Set teacher lesson data to cache
	if fetcher.caching {
		fetcher.cacheLock.Lock()
		fetcher.cache[teacher.ID] = &teacherLessons{teacher: teacher, lessons: lessons}
		fetcher.cacheLock.Unlock()
	}

	return teacher, lessons, nil
}

func (fetcher *LessonFetcher) parseTeacherAttribute(teacher *model.Teacher, rootNode *xmlpath.Node) {
	nameXPath := xmlpath.MustCompile(`./dt`)
	valueXPath := xmlpath.MustCompile(`./dd`)
	for iter := attributesXPath.Iter(rootNode); iter.Next(); {
		node := iter.Node()
		name, ok := nameXPath.String(node)
		if !ok {
			fetcher.logger.Error(
				fmt.Sprintf("Failed to parse teacher value: name=%v", name),
				zap.Uint("teacherID", uint(teacher.ID)),
			)
			continue
		}
		value, ok := valueXPath.String(node)
		if !ok {
			fetcher.logger.Error(
				fmt.Sprintf("Failed to parse teacher value: name=%v, value=%v", name, value),
				zap.Uint("teacherID", uint(teacher.ID)),
			)
			continue
		}
		if err := fetcher.setTeacherAttribute(teacher, strings.TrimSpace(name), strings.TrimSpace(value)); err != nil {
			fetcher.logger.Error(
				fmt.Sprintf("Failed to setTeacherAttribute: name=%v, value=%v", name, value),
				zap.Uint("teacherID", uint(teacher.ID)),
			)
		}
		//fmt.Printf("name = %v, value = %v\n", strings.TrimSpace(name), strings.TrimSpace(value))
	}
	//fmt.Printf("teacher = %+v\n", teacher)
}

func (fetcher *LessonFetcher) setTeacherAttribute(teacher *model.Teacher, name string, value string) error {
	switch name {
	case "国籍":
		c, found := fetcher.mCountries.GetByNameJA(value)
		if !found {
			return errors.NotFoundf("No MCountries for %v", value)
		}
		teacher.CountryID = c.ID
	case "誕生日":
		value = width.Narrow.String(value)
		if strings.TrimSpace(value) == "" {
			teacher.Birthday = time.Time{}
		} else {
			t, err := time.Parse("2006-01-02", value)
			if err != nil {
				return err
			}
			teacher.Birthday = t
		}
	case "性別":
		switch value {
		case "男性":
			teacher.Gender = "male" // TODO: enum
		case "女性":
			teacher.Gender = "female"
		default:
			return errors.Internalf("Unknown gender for %v", value)
		}
	case "経歴":
		yoe := -1
		switch value {
		case "1年未満":
			yoe = 0
		case "3年以上":
			yoe = 4
		default:
			value = strings.Replace(value, "年", "", -1)
			if v, err := strconv.ParseInt(width.Narrow.String(value), 10, 32); err == nil {
				yoe = int(v)
			} else {
				return errors.InternalWrapf(err, "Failed to convert to number: %v", value)
			}
		}
		teacher.YearsOfExperience = uint8(yoe)
	}
	return nil
}

func (fetcher *LessonFetcher) parseTeacherFavoriteCount(teacher *model.Teacher, rootNode *xmlpath.Node) {
	nameXPath := xmlpath.MustCompile(`//span[@id='fav_count']`)
	value, ok := nameXPath.String(rootNode)
	if !ok {
		fetcher.logger.Error(
			fmt.Sprintf("Failed to parse teacher favorite count"),
			zap.Uint("teacherID", uint(teacher.ID)),
		)
		return
	}
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		fetcher.logger.Error(
			fmt.Sprintf("Failed to parse teacher favorite count. It's not a number"),
			zap.Uint("teacherID", uint(teacher.ID)),
		)
		return
	}
	teacher.FavoriteCount = uint32(v)
}

func (fetcher *LessonFetcher) Close() {
	close(fetcher.semaphore)
}

func MustInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}

type MockTransport struct {
	sync.Mutex
	NumCalled int
	content   string
}

func NewMockTransport(path string) (*MockTransport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.Open failed: path=%v, err=%v", path, err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file failed: err=%v", err)
	}
	return &MockTransport{
		content: string(b),
	}, nil
}

func NewMockTransportFromHTML(content string) *MockTransport {
	return &MockTransport{
		content: content,
	}
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	t.NumCalled++
	t.Unlock()
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
		Status:     "200 OK",
	}
	resp.Header.Set("Content-Type", "text/html; charset=UTF-8")
	resp.Body = ioutil.NopCloser(strings.NewReader(t.content))
	return resp, nil
}
