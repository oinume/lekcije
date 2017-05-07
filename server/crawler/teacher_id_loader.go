package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

type teacherIDLoader interface {
	Load() ([]uint32, error)
}

type specificTeacherIDLoader struct {
	idString string
}

func (l *specificTeacherIDLoader) Load() ([]uint32, error) {
	sids := strings.Split(l.idString, ",")
	ids := make([]uint32, 0, len(sids))
	for _, id := range sids {
		i, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint32(i))
	}
	return ids, nil
}

type followedTeacherIDLoader struct {
	db *gorm.DB
}

func (l *followedTeacherIDLoader) Load() ([]uint32, error) {
	ids, err := model.NewFollowingTeacherService(l.db).FindTeacherIDs()
	if err != nil {
		return nil, err
	}
	return ids, nil
}

type scrapingOrder int

const (
	byRating scrapingOrder = iota + 1
	byNew
)

const (
	userAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html"
)

type scrapingTeacherIDLoader struct {
	order scrapingOrder
}

func (l *scrapingTeacherIDLoader) Load() ([]uint32, error) {
	u := "http://eikaiwa.dmm.com"
	u += "/list/?data%5Btab2%5D%5Bgender%5D=0&data%5Btab2%5D%5Bage%5D=%E5%B9%B4%E9%BD%A2&data%5Btab2%5D%5Bfree_word%5D=&tab=1&sort=4"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed to create HTTP request: url=%v", u)
	}
	req.Header.Set("User-Agent", userAgent)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed httpClient.Do(): url=%v", u)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Internalf("Unknown error in fetchContent: url=%v, status=%v", u, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", string(body))

	// TODO: Implement crawler

	return nil, nil
}
