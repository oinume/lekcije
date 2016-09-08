package notifier

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var lessonFetcher *fetcher.TeacherLessonFetcher

func init() {
	http.DefaultClient.Timeout = 5 * time.Second
	lessonFetcher = fetcher.NewTeacherLessonFetcher(http.DefaultClient, nil)
}

type Notifier struct {
	teachers map[uint32]*model.Teacher
}

func NewNotifier() *Notifier {
	return &Notifier{
		teachers: make(map[uint32]*model.Teacher, 1000),
	}
}

func (n *Notifier) SendNotification(user *model.User) error {
	teacherIds, err := model.FollowingTeacherService.FindTeacherIdsByUserId(user.Id)
	if err != nil {
		return errors.Wrapperf(err, "Failed to FindTeacherIdsByUserId(): userId=%v", user.Id)
	}

	availableLessonsPerTeacher := make(map[uint32][]*model.Lesson, 1000)
	for _, teacherId := range teacherIds {
		teacher, availableLessons, err := n.fetchAndExtractAvailableLessons(teacherId)
		if err != nil {
			return err
		}
		n.teachers[teacherId] = teacher
		for _, l := range availableLessons {
			fmt.Printf("available -> teacherId: %v, datetime:%v, status:%v \n", l.TeacherId, l.Datetime.Format("2006-01-02 15:04"), l.Status)
		}
		availableLessonsPerTeacher[teacherId] = availableLessons
	}

	// TODO: dry-run
	if err := n.sendNotificationToUser(user, availableLessonsPerTeacher); err != nil {
		return err
	}
	return nil
}

func (n *Notifier) fetchAndExtractAvailableLessons(teacherId uint32) (
	*model.Teacher, []*model.Lesson, error,
) {
	teacher, fetchedLessons, err := lessonFetcher.Fetch(teacherId)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Printf("--- %s ---\n", teacher.Name)
	//for _, lesson := range fetchedLessons {
	//	fmt.Printf("datetime = %v, status = %v\n", lesson.Datetime, lesson.Status)
	//}

	now := time.Now()
	fromDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.LocalTimezone())
	toDate := fromDate.Add(24 * 6 * time.Hour)
	lastFetchedLessons, err := model.LessonService.FindLessons(teacher.Id, fromDate, toDate)
	if err != nil {
		return nil, nil, err
	}
	availableLessons := model.LessonService.GetNewAvailableLessons(lastFetchedLessons, fetchedLessons)
	return teacher, availableLessons, nil

	//_, err = model.LessonService.UpdateLessons(fetchedLessons)
	//if err != nil {
	//	return err
	//}
}

func (n *Notifier) sendNotificationToUser(
	user *model.User,
	lessonsPerTeacher map[uint32][]*model.Lesson,
) error {
	lessonsCount := 0
	var teacherIds []int
	for teacherId, lessons := range lessonsPerTeacher {
		teacherIds = append(teacherIds, int(teacherId))
		lessonsCount += len(lessons)
	}
	if lessonsCount == 0 {
		// Don't send notification
		return nil
	}

	sort.Ints(teacherIds)
	var teacherIds2 []uint32
	var teacherNames []string
	for _, id := range teacherIds {
		teacherIds2 = append(teacherIds2, uint32(id))
		teacherNames = append(teacherNames, n.teachers[uint32(id)].Name)
	}

	t := template.New("email")
	t = template.Must(t.Parse(getEmailTemplate()))
	type TemplateData struct {
		TeacherIds        []uint32
		Teachers          map[uint32]*model.Teacher
		LessonsPerTeacher map[uint32][]*model.Lesson
	}
	data := &TemplateData{
		TeacherIds:        teacherIds2,
		Teachers:          n.teachers,
		LessonsPerTeacher: lessonsPerTeacher,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return errors.InternalWrapf(err, "Failed to execute template.")
	}
	//fmt.Printf("--- mail ---\n%s", body.String())

	subject := "[lekcije] Schedules of teacher " + strings.Join(teacherNames, ", ")
	sender := &EmailNotificationSender{}
	return sender.Send(user, subject, body.String())
}

func getEmailTemplate() string {
	return strings.TrimSpace(`
{{- range $teacherId := .TeacherIds }}
{{- $teacher := index $.Teachers $teacherId -}}
--- Teacher '{{ $teacher.Name }}' available lessons ---
  {{- $lessons := index $.LessonsPerTeacher $teacherId -}}
  {{- range $lesson := $lessons }}
{{ $lesson.Datetime }}
  {{- end }}

{{ end }}
	`)
}

type NotificationSender interface {
	Send(user *model.User, subject, body string) error
}

type EmailNotificationSender struct{}

func (s *EmailNotificationSender) Send(user *model.User, subject, body string) error {
	from := mail.NewEmail("noreply", "noreply@lampetty.net") // TODO: noreply@lekcije.com
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
		logger.AppLogger.Error(message)
		return errors.InternalWrapf(
			err,
			"Failed to send email by sendgrid: statusCode=%v",
			resp.StatusCode,
		)
	}

	return nil
}
