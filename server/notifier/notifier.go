package notifier

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/uber-go/zap"
)

type Notifier struct {
	db             *gorm.DB
	fetcher        *fetcher.TeacherLessonFetcher
	dryRun         bool
	lessonService  *model.LessonService
	teachers       map[uint32]*model.Teacher
	fetchedLessons map[uint32][]*model.Lesson
	sync.Mutex
}

func NewNotifier(db *gorm.DB, fetcher *fetcher.TeacherLessonFetcher, dryRun bool) *Notifier {
	return &Notifier{
		db:             db,
		fetcher:        fetcher,
		dryRun:         dryRun,
		teachers:       make(map[uint32]*model.Teacher, 1000),
		fetchedLessons: make(map[uint32][]*model.Lesson, 1000),
	}
}

func (n *Notifier) SendNotification(user *model.User) error {
	followingTeacherService := model.NewFollowingTeacherService(n.db)
	n.lessonService = model.NewLessonService(n.db)

	teacherIDs, err := followingTeacherService.FindTeacherIDsByUserID(user.ID)
	if err != nil {
		return errors.Wrapperf(err, "Failed to FindTeacherIDsByUserID(): userID=%v", user.ID)
	}
	if len(teacherIDs) != 0 {
		logger.App.Info(
			"Target teachers",
			zap.Uint("userID", uint(user.ID)),
			zap.String("teacherIDs", strings.Join(util.Uint32ToStringSlice(teacherIDs...), ",")),
		)
	}

	availableLessonsPerTeacher := make(map[uint32][]*model.Lesson, 1000)
	wg := &sync.WaitGroup{}
	for _, teacherID := range teacherIDs {
		wg.Add(1)
		go func(teacherID uint32) {
			defer wg.Done()
			teacher, fetchedLessons, newAvailableLessons, err := n.fetchAndExtractNewAvailableLessons(teacherID)
			if err != nil {
				switch err.(type) {
				case *errors.NotFound:
					// TODO: update teacher table flag
					// TODO: Not need to log
					logger.App.Warn("Cannot find teacher", zap.Uint("teacherID", uint(teacherID)))
				default:
					logger.App.Error("Cannot fetch teacher", zap.Uint("teacherID", uint(teacherID)), zap.Error(err))
				}
				return
			}

			n.Lock()
			defer n.Unlock()
			n.teachers[teacherID] = teacher
			if _, ok := n.fetchedLessons[teacherID]; !ok {
				n.fetchedLessons[teacherID] = make([]*model.Lesson, 0, 5000)
			}
			n.fetchedLessons[teacherID] = append(n.fetchedLessons[teacherID], fetchedLessons...)
			if len(newAvailableLessons) > 0 {
				availableLessonsPerTeacher[teacherID] = newAvailableLessons
			}
		}(teacherID)

		if err != nil {
			return err
		}
	}
	wg.Wait()

	if err := n.sendNotificationToUser(user, availableLessonsPerTeacher); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	return nil
}

// Returns teacher, fetchedLessons, newAvailableLessons, error
func (n *Notifier) fetchAndExtractNewAvailableLessons(teacherID uint32) (
	*model.Teacher, []*model.Lesson, []*model.Lesson, error,
) {
	teacher, fetchedLessons, err := n.fetcher.Fetch(teacherID)
	if err != nil {
		logger.App.Error(
			"fetcher.Fetch",
			zap.Uint("teacherID", uint(teacherID)), zap.Error(err),
		)
		return nil, nil, nil, err
	}
	logger.App.Debug(
		"fetcher.Fetch",
		zap.Uint("teacherID", uint(teacher.ID)),
		zap.Int("lessons", len(fetchedLessons)),
	)

	//fmt.Printf("fetchedLessons ---\n")
	//for _, l := range fetchedLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}

	now := time.Now()
	fromDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.LocalTimezone())
	toDate := fromDate.Add(24 * 6 * time.Hour)
	lastFetchedLessons, err := n.lessonService.FindLessons(teacher.ID, fromDate, toDate)
	if err != nil {
		return nil, nil, nil, err
	}
	//fmt.Printf("lastFetchedLessons ---\n")
	//for _, l := range lastFetchedLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}

	newAvailableLessons := n.lessonService.GetNewAvailableLessons(lastFetchedLessons, fetchedLessons)
	//fmt.Printf("newAvailableLessons ---\n")
	//for _, l := range newAvailableLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}
	return teacher, fetchedLessons, newAvailableLessons, nil
}

func (n *Notifier) sendNotificationToUser(
	user *model.User,
	lessonsPerTeacher map[uint32][]*model.Lesson,
) error {
	lessonsCount := 0
	var teacherIDs []int
	for teacherID, lessons := range lessonsPerTeacher {
		teacherIDs = append(teacherIDs, int(teacherID))
		lessonsCount += len(lessons)
	}
	if lessonsCount == 0 {
		// Don't send notification
		return nil
	}

	sort.Ints(teacherIDs)
	var teacherIDs2 []uint32
	var teacherNames []string
	for _, id := range teacherIDs {
		teacherIDs2 = append(teacherIDs2, uint32(id))
		teacherNames = append(teacherNames, n.teachers[uint32(id)].Name)
	}

	t := template.New("email")
	t = template.Must(t.Parse(getEmailTemplateJP()))
	type TemplateData struct {
		TeacherIDs        []uint32
		Teachers          map[uint32]*model.Teacher
		LessonsPerTeacher map[uint32][]*model.Lesson
		WebURL            string
	}
	data := &TemplateData{
		TeacherIDs:        teacherIDs2,
		Teachers:          n.teachers,
		LessonsPerTeacher: lessonsPerTeacher,
		WebURL:            config.WebURL(),
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return errors.InternalWrapf(err, "Failed to execute template.")
	}
	//fmt.Printf("--- mail ---\n%s", body.String())

	logger.App.Info("sendNotificationToUser", zap.String("email", user.Email.Raw()))
	//subject := "Schedule of teacher " + strings.Join(teacherNames, ", ")
	subject := strings.Join(teacherNames, ", ") + "の空きレッスンがあります"
	sender := &EmailNotificationSender{}
	return sender.Send(user, subject, body.String())
}

func getEmailTemplateJP() string {
	return strings.TrimSpace(`
{{- range $teacherID := .TeacherIDs }}
{{- $teacher := index $.Teachers $teacherID -}}
--- {{ $teacher.Name }} ---
  {{- $lessons := index $.LessonsPerTeacher $teacherID }}
  {{- range $lesson := $lessons }}
{{ $lesson.Datetime.Format "2006-01-02 15:04" }}
  {{- end }}

レッスンの予約はこちらから:
<a href="http://eikaiwa.dmm.com/teacher/index/{{ $teacherID }}/">PC</a>
<a href="http://eikaiwa.dmm.com/teacher/schedule/{{ $teacherID }}/">Mobile</a>

{{ end }}
空きレッスンの通知の解除は<a href="{{ .WebURL }}/me">こちら</a>
	`)
}

func getEmailTemplateEN() string {
	return strings.TrimSpace(`
{{- range $teacherID := .TeacherIDs }}
{{- $teacher := index $.Teachers $teacherID -}}
--- {{ $teacher.Name }} ---
  {{- $lessons := index $.LessonsPerTeacher $teacherID }}
  {{- range $lesson := $lessons }}
{{ $lesson.Datetime.Format "2006-01-02 15:04" }}
  {{- end }}

Reserve here:
<a href="http://eikaiwa.dmm.com/teacher/index/{{ $teacherID }}/">PC</a>
<a href="http://eikaiwa.dmm.com/teacher/schedule/{{ $teacherID }}/">Mobile</a>
{{ end }}
Click <a href="{{ .WebURL }}/me">here</a> if you want to stop notification of the teacher.
	`)
}

func (n *Notifier) Close() {
	defer n.fetcher.Close()
	defer func() {
		if n.dryRun {
			return
		}
		for teacherID, lessons := range n.fetchedLessons {
			if _, err := n.lessonService.UpdateLessons(lessons); err != nil {
				logger.App.Error(
					"An error ocurred in Notifier.Close",
					zap.Error(err), zap.Uint("teacherID", uint(teacherID)),
				)
			}
		}
	}()
}

type NotificationSender interface {
	Send(user *model.User, subject, body string) error
}

type EmailNotificationSender struct{}

func (s *EmailNotificationSender) Send(user *model.User, subject, body string) error {
	from := mail.NewEmail("lekcije", "lekcije@lekcije.com")
	to := mail.NewEmail(user.Name, user.Email.Raw())
	content := mail.NewContent("text/html", strings.Replace(body, "\n", "<br>", -1))
	m := mail.NewV3MailInit(from, subject, to, content)

	req := sendgrid.GetRequest(
		os.Getenv("SENDGRID_API_KEY"),
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	req.Method = "POST"
	req.Body = mail.GetRequestBody(m)
	resp, err := sendgrid.API(req)
	if err != nil {
		return errors.InternalWrapf(err, "Failed to send email by sendgrid")
	}
	if resp.StatusCode >= 300 {
		message := fmt.Sprintf(
			"Failed to send email by sendgrid: statusCode=%v, body=%v",
			resp.StatusCode, strings.Replace(resp.Body, "\n", "\\n", -1),
		)
		logger.App.Error(message)
		return errors.InternalWrapf(
			err,
			"Failed to send email by sendgrid: statusCode=%v",
			resp.StatusCode,
		)
	}

	return nil
}
