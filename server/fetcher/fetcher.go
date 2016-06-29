package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/model"
	"gopkg.in/xmlpath.v2"
)

const (
	urlBase   = "http://eikaiwa.dmm.com/teacher/index/%v/"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/601.6.17 (KHTML, like Gecko) Version/9.1.1 Safari/601.6.17"
)

var (
	jst            = time.FixedZone("Asia/Tokyo", 9*60*60)
	titleXPath     = xmlpath.MustCompile(`//title`)
	lessonXPath    = xmlpath.MustCompile("//ul[@class='oneday']//li")
	classAttrXPath = xmlpath.MustCompile("@class")
)

type TeacherLessonFetcher struct {
	httpClient *http.Client
}

func NewTeacherLessonFetcher(httpClient *http.Client) *TeacherLessonFetcher {
	client := httpClient
	if client == nil {
		client = http.DefaultClient
		client.Timeout = 5 * time.Second
		// TODO: retry
	}
	return &TeacherLessonFetcher{
		httpClient: client,
	}
}

func (fetcher *TeacherLessonFetcher) Fetch(teacherId uint32) (*model.Teacher, []*model.Lesson, error) {
	targetUrl := fmt.Sprintf(urlBase, teacherId)
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := fetcher.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO: pkg/errors
		return nil, nil, fmt.Errorf("fetch error: url=%v, status=%v", targetUrl, resp.StatusCode)
	}
	return fetcher.parseHtml(teacherId, resp.Body)
}

func (fetcher *TeacherLessonFetcher) parseHtml(teacherId uint32, html io.Reader) (*model.Teacher, []*model.Lesson, error) {
	root, err := xmlpath.ParseHTML(html)
	if err != nil {
		return nil, nil, err
	}
	teacher := &model.Teacher{
		Id: teacherId,
	}

	// teacher name
	if title, ok := titleXPath.String(root); ok {
		teacher.Name = strings.Trim(strings.Split(title, "-")[0], " ")
	} else {
		return nil, nil, fmt.Errorf("failed to fetch teacher's name: url=%v", fmt.Sprintf(urlBase, teacherId))
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
		fmt.Printf("text = '%v', timeClass = '%v'\n", text, timeClass)

		// blank, reservable, reserved
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
			status := ""
			switch text { // TODO: enum
			case "終了":
				status = "finished"
			case "予約済", "休講":
				status = "reserved"
			case "予約可":
				status = "reservable"
			default:
				// TODO: error
			}
			fmt.Printf("dt = %v, status=%v\n", dt, status) // TODO: logging

			lessons = append(lessons, &model.Lesson{
				TeacherId: teacher.Id,
				Datetime:  dt,
				Status:    status,
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

/*
       # schedule, reservation
       original_date = datetime.date.today()
       date = copy.copy(original_date)
       time_items = root.xpath("//ul[@class='oneday']//li")
       schedules = []
       logger.debug("--- teacher id={}, name={} ---".format(teacher_id, name))
       for time_item in time_items:
           time_class = time_item.attrib["class"]
           text = time_item.text_content().strip()
           # logger.debug("web {time_class}:{text}".format(**locals()))
           # blank, reservable, reserved
           if time_class == "date":
               match = re.match(r"([\d]+)月([\d]+)日(.+)", text)
               if match:
                   original_date = date.replace(date.year, int(match.group(1)), int(match.group(2)))
                   date = copy.copy(original_date)
           elif time_class.startswith("t-") and text != "":
               tmp = time_class.split("-")
               hour, minute = int(tmp[1]), int(tmp[2])
               if hour >= 24:
                   # 24:30 -> 00:30
                   hour -= 24
                   if date == original_date:
                       # Set date to next day for 24:30
                       date += datetime.timedelta(days=1)
               dt = datetime.datetime(date.year, date.month, date.day, hour, minute, 0, 0)
               if text == "終了":
                   status = "finished"
               elif text in ("予約済", "休講"):  # TODO: Add this status to enum
                   status = "reserved"
               elif text == "予約可":
                   status = "reservable"
               else:
                   raise(Exception("Unknown schedule text:{}".format(text)))

               logger.debug("dt={dt}, status={status}".format(**locals()))
               schedule = Schedule(teacher.id, dt, ScheduleStatus[status])
               schedules.append(schedule)
           else:
               pass
       return teacher, schedules

   def close(self):
       self._session.close()

*/
