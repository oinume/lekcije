package notifier

import (
	"fmt"
	"net/http"
	"sort"
	"text/template"
	"time"

	"bytes"
	"strings"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/model"
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
	var teacherIds []int
	for teacherId := range lessonsPerTeacher {
		teacherIds = append(teacherIds, int(teacherId))
	}
	fmt.Printf("teacherIds = %+v\n", teacherIds)
	sort.Ints(teacherIds)
	var teacherIds2 []uint32
	for _, id := range teacherIds {
		println("id = ", id)
		teacherIds2 = append(teacherIds2, uint32(id))
	}
	fmt.Printf("teacherIds2 = %+v\n", teacherIds2)

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

	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return errors.InternalWrapf(err, "Failed to execute template.")
	}
	fmt.Printf("--- mail ---\n%s", b.String())
	return nil
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
