package dmm_eikaiwa

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Songmu/retry"
	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/v4/types"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"golang.org/x/text/width"
	"gopkg.in/xmlpath.v2"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

const (
	userAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html"
)

var (
	_                 = fmt.Print
	defaultHTTPClient = &http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: redirectErrorFunc,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			Proxy:               http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 1200 * time.Second,
			}).DialContext,
			IdleConnTimeout:     1200 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				ClientSessionCache: tls.NewLRUClientSessionCache(100),
			},
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	redirectErrorFunc = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	teacherNameXPath = xmlpath.MustCompile(`//div[@class='area-detail']/h1/text()`)
	attributesXPath  = xmlpath.MustCompile(`//div[@class='confirm low']/dl`)
	lessonXPath      = xmlpath.MustCompile(`//ul[@class='oneday']//li`)
	classAttrXPath   = xmlpath.MustCompile(`@class`)
	ratingXPath      = xmlpath.MustCompile(`//p[@class='ui-star-rating-text']/strong/text()`)
	reviewCountXPath = xmlpath.MustCompile(`//p[@class='ui-star-rating-text']/text()`)
	newTeacherXPath  = xmlpath.MustCompile(`//div[@class='favorite-list-box-wrap']/img[@class='new_teacher']`)
)

type teacherLessons struct {
	teacher *model2.Teacher
	lessons []*model2.Lesson
}

type lessonFetcher struct {
	httpClient   *http.Client
	semaphore    chan struct{}
	caching      bool
	cache        map[uint]*teacherLessons
	cacheLock    *sync.RWMutex
	logger       *zap.Logger
	mCountryList *model2.MCountryList
}

func NewLessonFetcher(
	httpClient *http.Client,
	concurrency int,
	caching bool,
	mCountryList *model2.MCountryList,
	log *zap.Logger,
) repository.LessonFetcher {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	return &lessonFetcher{
		httpClient:   httpClient,
		semaphore:    make(chan struct{}, concurrency),
		caching:      caching,
		cache:        make(map[uint]*teacherLessons, 5000),
		cacheLock:    new(sync.RWMutex),
		mCountryList: mCountryList,
		logger:       log,
	}
}

func (f *lessonFetcher) Fetch(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error) {
	ctx, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "LessonFetcher.Fetch")
	span.SetAttributes(attribute.KeyValue{
		Key:   "teacherID",
		Value: attribute.Int64Value(int64(teacherID)),
	})
	defer span.End()

	f.semaphore <- struct{}{}
	defer func() {
		<-f.semaphore
	}()

	// Check cache
	if f.caching {
		f.cacheLock.RLock()
		if c, ok := f.cache[teacherID]; ok {
			f.cacheLock.RUnlock()
			return c.teacher, c.lessons, nil
		}
		f.cacheLock.RUnlock()
	}

	teacher := model2.NewTeacher(teacherID)
	var content io.ReadCloser
	err := retry.Retry(2, 300*time.Millisecond, func() error {
		var err error
		content, err = f.fetchContent(ctx, teacher.URL())
		return err
	})
	defer content.Close()
	if err != nil {
		return nil, nil, err
	}

	_, lessons, err := f.parseHTML(teacher, content)
	if err != nil {
		return nil, nil, err
	}
	if len(lessons) > 0 {
		teacher.LastLessonAt = lessons[len(lessons)-1].Datetime
	}
	return teacher, lessons, nil
}

func (f *lessonFetcher) fetchContent(ctx context.Context, url string) (io.ReadCloser, error) {
	clientTrace := otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutSubSpans())
	ctx = httptrace.WithClientTrace(ctx, clientTrace)
	nopCloser := io.NopCloser(strings.NewReader(""))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nopCloser, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to create HTTP request: url=%v", url),
		)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nopCloser, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed httpClient.Do(): url=%v", url),
		)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.Body, nil
	case http.StatusMovedPermanently, http.StatusFound:
		_ = resp.Body.Close()
		return nopCloser, errors.NewNotFoundError(
			errors.WithMessagef("Teacher not found: url=%v, statusCode=%v", url, resp.StatusCode),
		)
	default:
		_ = resp.Body.Close()
		return nopCloser, errors.NewInternalError(
			errors.WithMessagef(
				"Unknown error in fetchContent: url=%v, statusCode=%v, status=%v",
				url, resp.StatusCode, resp.Status,
			),
		)
	}
}

func (f *lessonFetcher) parseHTML(
	teacher *model2.Teacher,
	html io.Reader,
) (*model2.Teacher, []*model2.Lesson, error) {
	root, err := xmlpath.ParseHTML(html)
	if err != nil {
		return nil, nil, err
	}

	// teacher name
	if teacherName, ok := teacherNameXPath.String(root); ok {
		teacher.Name = teacherName
	} else {
		return nil, nil, fmt.Errorf("failed to fetch teacher's name: url=%v", teacher.URL())
	}

	// Nationality, birthday, etc...
	f.parseTeacherAttribute(teacher, root)
	if !teacher.IsJapanese() { // Japanese teachers don't have favorite count
		// FavoriteCount
		f.parseTeacherFavoriteCount(teacher, root)
	}
	// Rating
	f.parseTeacherRating(teacher, root)
	// ReviewCount
	f.parseTeacherReviewCount(teacher, root)

	dateRegexp := regexp.MustCompile(`([\d]+)月([\d]+)日(.+)`)
	lessons := make([]*model2.Lesson, 0, 1000)
	now := time.Now().In(config.LocalLocation())
	originalDate := time.Now().In(config.LocalLocation()).Truncate(24 * time.Hour)
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
		f.logger.Debug("Scraping as", zap.String("timeClass", timeClass), zap.String("text", text))

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
					config.LocalLocation(),
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
				config.LocalLocation(),
			)
			status := model2.LessonStatuses.MustValueForAlias(text)
			f.logger.Debug(
				"lesson",
				zap.String("dt", dt.Format("2006-01-02 15:04")),
				zap.String("status", model.LessonStatuses.MustName(status)),
			)
			lessons = append(lessons, &model2.Lesson{
				TeacherID: teacher.ID,
				Datetime:  dt,
				Status:    model2.LessonStatuses.MustName(status),
			})
		}
		// TODO: else
	}

	// Set teacher lesson data to cache
	if f.caching {
		f.cacheLock.Lock()
		f.cache[teacher.ID] = &teacherLessons{teacher: teacher, lessons: lessons}
		f.cacheLock.Unlock()
	}

	return teacher, lessons, nil
}

func (f *lessonFetcher) parseTeacherAttribute(teacher *model2.Teacher, rootNode *xmlpath.Node) {
	nameXPath := xmlpath.MustCompile(`./dt`)
	valueXPath := xmlpath.MustCompile(`./dd`)
	for iter := attributesXPath.Iter(rootNode); iter.Next(); {
		node := iter.Node()
		name, ok := nameXPath.String(node)
		if !ok {
			f.logger.Error(
				fmt.Sprintf("Failed to parse teacher value: name=%v", name),
				zap.Uint("teacherID", teacher.ID),
			)
			continue
		}
		value, ok := valueXPath.String(node)
		if !ok {
			f.logger.Error(
				fmt.Sprintf("Failed to parse teacher value: name=%v, value=%v", name, value),
				zap.Uint("teacherID", teacher.ID),
			)
			continue
		}
		if err := f.setTeacherAttribute(teacher, strings.TrimSpace(name), strings.TrimSpace(value)); err != nil {
			f.logger.Error(
				fmt.Sprintf("Failed to setTeacherAttribute: name=%v, value=%v", name, value),
				zap.Uint("teacherID", teacher.ID),
			)
		}
		//fmt.Printf("name = %v, value = %v\n", strings.TrimSpace(name), strings.TrimSpace(value))
	}
	//fmt.Printf("teacher = %+v\n", teacher)
}

func (f *lessonFetcher) setTeacherAttribute(teacher *model2.Teacher, name string, value string) error {
	switch name {
	case "国籍":
		c, found := f.mCountryList.GetByNameJA(value)
		if !found {
			return errors.NewNotFoundError(errors.WithMessage(fmt.Sprintf("No MCountries for %v", value)))
		}
		teacher.CountryID = int16(c.ID) // TODO: teacher.CountryID must be uint16
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
			return errors.NewInternalError(
				errors.WithMessagef("Unknown gender for %v", value),
			)
		}
	case "経歴":
		var yoe int
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
				return errors.NewInternalError(
					errors.WithError(err),
					errors.WithMessagef("Failed to convert to number: %v", value),
				)
			}
		}
		teacher.YearsOfExperience = int8(yoe) // TODO: teacher.YearsOfExperience must be uint8
	}
	return nil
}

func (f *lessonFetcher) parseTeacherFavoriteCount(teacher *model2.Teacher, rootNode *xmlpath.Node) {
	favCountXPath := xmlpath.MustCompile(`//span[@id='fav_count']`)
	value, ok := favCountXPath.String(rootNode)
	if !ok {
		f.logger.Error(
			"Failed to parse teacher favorite count",
			zap.Uint("teacherID", uint(teacher.ID)),
		)
		return
	}
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		f.logger.Error(
			"Failed to parse teacher favorite count. It's not a number",
			zap.Uint("teacherID", uint(teacher.ID)),
		)
		return
	}
	teacher.FavoriteCount = uint(v)
}

func (f *lessonFetcher) parseTeacherRating(teacher *model2.Teacher, rootNode *xmlpath.Node) {
	value, ok := ratingXPath.String(rootNode)
	if !ok {
		if _, ok := newTeacherXPath.String(rootNode); !ok {
			f.logger.Error(
				"Failed to parse teacher rating",
				zap.Uint("teacherID", teacher.ID),
				zap.String("value", value),
			)
		}
		// Give up to obtain rating
		return
	}
	rating, err := strconv.ParseFloat(value, 32)
	if err != nil {
		f.logger.Error(
			"Failed to parse teacher rating. It's not a number",
			zap.Uint("teacherID", teacher.ID),
		)
		return
	}
	teacher.Rating = types.NullDecimal{Big: decimal.New(int64(rating*100), 2)}
}

func (f *lessonFetcher) parseTeacherReviewCount(teacher *model2.Teacher, rootNode *xmlpath.Node) {
	value, ok := reviewCountXPath.String(rootNode)
	if !ok {
		if _, ok := newTeacherXPath.String(rootNode); !ok {
			f.logger.Error(
				"Failed to parse teacher review count",
				zap.Uint("teacherID", teacher.ID),
				zap.String("value", value),
			)
		}
		// Give up to obtain rating
		return
	}
	value = strings.TrimPrefix(value, "(")
	value = strings.TrimSuffix(value, ")")
	reviewCount, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		f.logger.Error(
			"Failed to parse teacher review count. It's not a number",
			zap.Uint("teacherID", teacher.ID),
			zap.String("value", value),
		)
		return
	}
	teacher.ReviewCount = uint(reviewCount)
}

func (f *lessonFetcher) Close() {
	close(f.semaphore)
}

func MustInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}
