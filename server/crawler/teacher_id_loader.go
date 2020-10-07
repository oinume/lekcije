package crawler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"gopkg.in/xmlpath.v2"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

type TeacherIDLoader interface {
	Load(cursor string) ([]uint32, string, error)
	GetInitialCursor() string
}

func NewSpecificTeacherIDLoader(idString string) TeacherIDLoader {
	return &specificTeacherIDLoader{
		idString: idString,
	}
}

type specificTeacherIDLoader struct {
	idString string
}

func (l *specificTeacherIDLoader) GetInitialCursor() string {
	return "dummy"
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

func NewFollowedTeacherIDLoader(db *gorm.DB) TeacherIDLoader {
	return &followedTeacherIDLoader{db: db}
}

type followedTeacherIDLoader struct {
	db *gorm.DB
}

func (l *followedTeacherIDLoader) GetInitialCursor() string {
	return "dummy"
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
	ByRating scrapingOrder = iota + 1
	ByNew
)

const (
	userAgent         = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html"
	scrapingURLPrefix = "http://eikaiwa.dmm.com/list/?"
)

var (
	defaultScrapingHTTPClient = &http.Client{}
	profileXPath              = xmlpath.MustCompile(`//li[@class='profile']`)
	idXPath                   = xmlpath.MustCompile(`@id`)
	paginationXPath           = xmlpath.MustCompile(`//div[@class='list-boxpagenation']/ul/li`)
	hrefXPath                 = xmlpath.MustCompile(`a/@href`)
	currentPageXpath          = xmlpath.MustCompile(`span`)
)

func NewScrapingTeacherIDLoader(order scrapingOrder, httpClient *http.Client) *scrapingTeacherIDLoader {
	if httpClient == nil {
		httpClient = defaultScrapingHTTPClient
	}
	return &scrapingTeacherIDLoader{
		order:      order,
		httpClient: httpClient,
	}
}

type scrapingTeacherIDLoader struct {
	order      scrapingOrder
	httpClient *http.Client
}

func (l *scrapingTeacherIDLoader) GetInitialCursor() string {
	return "data%5Btab2%5D%5Bgender%5D=0&data%5Btab2%5D%5Bage%5D=%E5%B9%B4%E9%BD%A2&data%5Btab2%5D%5Bfree_word%5D=&tab=1&sort=4"
}

// https://eikaiwa.dmm.com/list/?<cursor>
func (l *scrapingTeacherIDLoader) Load(cursor string) ([]uint32, string, error) {
	url := scrapingURLPrefix + cursor
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to create HTTP request: url=%v", url),
		)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, "", errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed httpClient.Do(): url=%v", url),
		)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.NewInternalError(
			errors.WithMessagef("Unknown error in fetchContent: url=%v, status=%v", url, resp.StatusCode),
		)
	}
	defer resp.Body.Close()

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return nil, "", errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to parse HTML: url=%v", url),
		)
	}

	ids := make([]uint32, 0, 100)
	for iter := profileXPath.Iter(root); iter.Next(); {
		node := iter.Node()
		id, ok := idXPath.String(node)
		if !ok {
			continue
		}
		v, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return nil, "", errors.NewInternalError(
				errors.WithMessagef("Failed to parse id: id=%v", id),
			)
		}
		ids = append(ids, uint32(v))
	}

	cursor = l.getNextCursor(root)
	return ids, cursor, nil
}

func (l *scrapingTeacherIDLoader) getNextCursor(root *xmlpath.Node) string {
	/*
		<div>
		<ul>
		  <li><a href="...">1</a></li>
		  <li><a href="...">2</a></li>
		  <li><p class="dot">...</p></li>
		  <li><span>4<span></li>           <-- current page
		  <li><a href="...">5</ap></li>    <-- want to get this a.href!
		</ul>
	*/
	currentPagePassed := false
	for iter := paginationXPath.Iter(root); iter.Next(); {
		node := iter.Node()
		//fmt.Printf("node = %v\n", node.String())
		if _, ok := currentPageXpath.String(node); ok {
			//fmt.Printf("currentPage = %v\n", n)
			currentPagePassed = true
		}
		if href, ok := hrefXPath.String(node); ok {
			if !currentPagePassed {
				continue
			}
			//fmt.Printf("href = %v\n", href)
			return strings.Replace(href, "/list/?", "", -1)
		}
	}

	return ""
}
