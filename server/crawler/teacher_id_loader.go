package crawler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"gopkg.in/xmlpath.v2"
)

type teacherIDLoader interface {
	Load(cursor string) ([]uint32, string, error)
}

type specificTeacherIDLoader struct {
	idString string
}

func (l *specificTeacherIDLoader) Load(cursor string) ([]uint32, string, error) {
	sids := strings.Split(l.idString, ",")
	ids := make([]uint32, 0, len(sids))
	for _, id := range sids {
		i, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return nil, "", err
		}
		ids = append(ids, uint32(i))
	}
	return ids, "", nil
}

type followedTeacherIDLoader struct {
	db *gorm.DB
}

func (l *followedTeacherIDLoader) Load(cursor string) ([]uint32, string, error) {
	ids, err := model.NewFollowingTeacherService(l.db).FindTeacherIDs()
	if err != nil {
		return nil, "", err
	}
	return ids, "", nil // TODO: cursor
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

func (l *scrapingTeacherIDLoader) Load(cursor string) ([]uint32, string, error) {
	u := "http://eikaiwa.dmm.com"
	u += "/list/?data%5Btab2%5D%5Bgender%5D=0&data%5Btab2%5D%5Bage%5D=%E5%B9%B4%E9%BD%A2&data%5Btab2%5D%5Bfree_word%5D=&tab=1&sort=4"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, "", errors.InternalWrapf(err, "Failed to create HTTP request: url=%v", u)
	}
	req.Header.Set("User-Agent", userAgent)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", errors.InternalWrapf(err, "Failed httpClient.Do(): url=%v", u)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.Internalf("Unknown error in fetchContent: url=%v, status=%v", u, resp.StatusCode)
	}
	defer resp.Body.Close()

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return nil, "", errors.Internalf("Failed to parse HTML: url=%v", u)
	}

	profileXPath := xmlpath.MustCompile(`//li[@class='profile']`)
	idXPath := xmlpath.MustCompile(`@id`)

	ids := make([]uint32, 0, 100)
	for iter := profileXPath.Iter(root); iter.Next(); {
		node := iter.Node()
		id, ok := idXPath.String(node)
		if !ok {
			continue
		}
		v, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return nil, "", errors.Internalf("Failed to parse id: id=%v", id)
		}
		ids = append(ids, uint32(v))
	}

	// TODO: scrape next page cursor

	return ids, "", nil
}
