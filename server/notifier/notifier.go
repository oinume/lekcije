package notifier

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/emailer"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/stopwatch"
	"github.com/oinume/lekcije/server/util"
	"go.uber.org/zap"
)

type Notifier struct {
	db              *gorm.DB
	fetcher         *fetcher.LessonFetcher
	dryRun          bool
	lessonService   *model.LessonService
	teachers        map[uint32]*model.Teacher
	fetchedLessons  map[uint32][]*model.Lesson
	sender          emailer.Sender
	senderWaitGroup *sync.WaitGroup
	stopwatch       stopwatch.Stopwatch
	sync.Mutex
}

type teachersAndLessons struct {
	data         map[uint32]*model.TeacherLessons
	lessonsCount int
	teacherIDs   []uint32
}

func (tal *teachersAndLessons) CountLessons() int {
	count := 0
	for _, l := range tal.data {
		count += len(l.Lessons)
	}
	return count
}

// Filter out by NotificationTimeSpanList.
// If a lesson is within NotificationTimeSpanList, it'll be included in returned value.
func (tal *teachersAndLessons) FilterBy(list model.NotificationTimeSpanList) *teachersAndLessons {
	if len(list) == 0 {
		return tal
	}
	ret := NewTeachersAndLessons(len(tal.data))
	for teacherID, tl := range tal.data {
		lessons := make([]*model.Lesson, 0, len(tl.Lessons))
		for _, lesson := range tl.Lessons {
			dt := lesson.Datetime
			t, _ := time.Parse("15:04", fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute()))
			if list.Within(t) {
				lessons = append(lessons, lesson)
			}
		}
		ret.data[teacherID] = model.NewTeacherLessons(tl.Teacher, lessons)
	}
	return ret
}

func (tal *teachersAndLessons) String() string {
	b := new(bytes.Buffer)
	for _, tl := range tal.data {
		fmt.Fprintf(b, "Teacher: %+v", tl.Teacher)
		fmt.Fprint(b, ", Lessons:")
		for _, l := range tl.Lessons {
			fmt.Fprintf(b, " {%+v}", l)
		}
	}
	return b.String()
}

func NewTeachersAndLessons(length int) *teachersAndLessons {
	return &teachersAndLessons{
		data:         make(map[uint32]*model.TeacherLessons, length),
		lessonsCount: -1,
		teacherIDs:   make([]uint32, 0, length),
	}
}

func NewNotifier(db *gorm.DB, fetcher *fetcher.LessonFetcher, dryRun bool, sender emailer.Sender) *Notifier {
	return &Notifier{
		db:              db,
		fetcher:         fetcher,
		dryRun:          dryRun,
		teachers:        make(map[uint32]*model.Teacher, 1000),
		fetchedLessons:  make(map[uint32][]*model.Lesson, 1000),
		sender:          sender,
		senderWaitGroup: &sync.WaitGroup{},
		stopwatch:       stopwatch.NewSync().Start(),
	}
}

func (n *Notifier) SendNotification(user *model.User) error {
	followingTeacherService := model.NewFollowingTeacherService(n.db)
	n.lessonService = model.NewLessonService(n.db)
	const maxFetchErrorCount = 5
	teacherIDs, err := followingTeacherService.FindTeacherIDsByUserID(user.ID, maxFetchErrorCount)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to FindTeacherIDsByUserID(): userID=%v", user.ID),
		)
	}
	n.stopwatch.Mark(fmt.Sprintf("FindTeacherIDsByUserID:%d", user.ID))

	if len(teacherIDs) == 0 {
		return nil
	}

	// Comment out due to papertrail limit
	//logger.App.Info("n", zap.Uint("userID", uint(user.ID)), zap.Int("teachers", len(teacherIDs)))

	//availableTeachersAndLessons := make(map[uint32][]*model.Lesson, 1000)
	availableTeachersAndLessons := NewTeachersAndLessons(1000)
	wg := &sync.WaitGroup{}
	for _, teacherID := range teacherIDs {
		wg.Add(1)
		go func(teacherID uint32) {
			defer n.stopwatch.Mark(fmt.Sprintf("fetchAndExtractNewAvailableLessons:%d", teacherID))
			defer wg.Done()
			fetched, newAvailable, err := n.fetchAndExtractNewAvailableLessons(teacherID)
			if err != nil {
				if errors.IsNotFound(err) {
					if err := model.NewTeacherService(n.db).IncrementFetchErrorCount(teacherID, 1); err != nil {
						logger.App.Error(
							"IncrementFetchErrorCount failed",
							zap.Uint("teacherID", uint(teacherID)), zap.Error(err),
						)
					}
					logger.App.Warn("Cannot find teacher", zap.Uint("teacherID", uint(teacherID)))
				}
				// TODO: Handle a case eikaiwa.dmm.com is down
				logger.App.Error("Cannot fetch teacher", zap.Uint("teacherID", uint(teacherID)), zap.Error(err))
				return
			}

			n.Lock()
			defer n.Unlock()
			n.teachers[teacherID] = fetched.Teacher
			if _, ok := n.fetchedLessons[teacherID]; !ok {
				n.fetchedLessons[teacherID] = make([]*model.Lesson, 0, 5000)
			}
			n.fetchedLessons[teacherID] = append(n.fetchedLessons[teacherID], fetched.Lessons...)
			if len(newAvailable.Lessons) > 0 {
				availableTeachersAndLessons.data[teacherID] = newAvailable
			}
			//fmt.Printf("go routine finished: user=%v\n", user.ID)
		}(teacherID)

		if err != nil {
			return err
		}
	}
	wg.Wait()

	notificationTimeSpanService := model.NewNotificationTimeSpanService(n.db)
	timeSpans, err := notificationTimeSpanService.FindByUserID(user.ID)
	if err != nil {
		return err
	}
	filteredAvailable := availableTeachersAndLessons.FilterBy(model.NotificationTimeSpanList(timeSpans))
	if err := n.sendNotificationToUser(user, filteredAvailable); err != nil {
		return err
	}

	time.Sleep(150 * time.Millisecond)
	n.stopwatch.Mark("sleep")

	return nil
}

// Returns teacher, fetchedLessons, newAvailableLessons, error
func (n *Notifier) fetchAndExtractNewAvailableLessons(teacherID uint32) (
	*model.TeacherLessons, *model.TeacherLessons, error,
) {
	teacher, fetchedLessons, err := n.fetcher.Fetch(teacherID)
	if err != nil {
		return nil, nil, err
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

	now := time.Now().In(config.LocalTimezone())
	fromDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.LocalTimezone())
	toDate := fromDate.Add(24 * 6 * time.Hour)
	lastFetchedLessons, err := n.lessonService.FindLessons(teacher.ID, fromDate, toDate)
	if err != nil {
		return nil, nil, err
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
	return model.NewTeacherLessons(teacher, fetchedLessons),
		model.NewTeacherLessons(teacher, newAvailableLessons),
		nil
	//return teacher, fetchedLessons, newAvailableLessons, nil
}

func (n *Notifier) sendNotificationToUser(
	user *model.User, lessonsPerTeacher *teachersAndLessons,
) error {
	lessonsCount := 0
	var teacherIDs []int
	for teacherID, l := range lessonsPerTeacher.data {
		teacherIDs = append(teacherIDs, int(teacherID))
		lessonsCount += len(l.Lessons)
	}
	if lessonsPerTeacher.CountLessons() == 0 {
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

	// TODO: getEmailTemplate as a static file
	t := emailer.NewTemplate("notifier", getEmailTemplateJP())
	data := struct {
		To                string
		TeacherNames      string
		TeacherIDs        []uint32
		Teachers          map[uint32]*model.Teacher
		LessonsPerTeacher map[uint32]*model.TeacherLessons
		WebURL            string
	}{
		To:                user.Email,
		TeacherNames:      strings.Join(teacherNames, ", "),
		TeacherIDs:        teacherIDs2,
		Teachers:          n.teachers,
		LessonsPerTeacher: lessonsPerTeacher.data,
		WebURL:            config.WebURL(),
	}
	email, err := emailer.NewEmailFromTemplate(t, data)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to create emailer.Email from template: to=%v", user.Email),
		)
	}
	email.SetCustomArg("email_type", model.EmailTypeNewLessonNotifier)
	email.SetCustomArg("user_id", fmt.Sprint(user.ID))
	email.SetCustomArg("teacher_ids", strings.Join(util.Uint32ToStringSlice(teacherIDs2...), ","))
	//fmt.Printf("--- mail ---\n%s", email.BodyString())
	n.stopwatch.Mark("emailer.NewEmailFromTemplate")

	logger.App.Info("sendNotificationToUser", zap.String("email", user.Email))

	n.senderWaitGroup.Add(1)
	go func(email *emailer.Email) {
		defer n.stopwatch.Mark(fmt.Sprintf("sender.Send:%d", user.ID))
		defer n.senderWaitGroup.Done()
		if err := n.sender.Send(email); err != nil {
			logger.App.Error(
				"Failed to sendNotificationToUser",
				zap.String("email", user.Email), zap.Error(err),
			)
		}
	}(email)

	return nil
}

func getEmailTemplateJP() string {
	return strings.TrimSpace(`
From: lekcije <lekcije@lekcije.com>
To: {{ .To }}
Subject: {{ .TeacherNames }}の空きレッスンがあります
Body: text/html
{{ range $teacherID := .TeacherIDs }}
{{- $teacher := index $.Teachers $teacherID -}}
--- {{ $teacher.Name }} ---
  {{- $tal := index $.LessonsPerTeacher $teacherID }}
  {{- range $lesson := $tal.Lessons }}
{{ $lesson.Datetime.Format "2006-01-02 15:04" }}
  {{- end }}

レッスンの予約はこちらから:
<a href="http://eikaiwa.dmm.com/teacher/index/{{ $teacherID }}/">PC</a>
<a href="http://eikaiwa.dmm.com/teacher/schedule/{{ $teacherID }}/">Mobile</a>

{{ end }}
PR ────────────────
<a href="https://lekcije.amebaownd.com/posts/3780950" target="_blank" rel="nofollow">設定画面からレッスン希望時間帯が設定できるようになりました</a>
PR ────────────────

空きレッスンの通知の解除は<a href="{{ .WebURL }}/me">こちら</a>

<a href="https://goo.gl/forms/CIGO3kpiQCGjtFD42">お問い合わせ</a>
	`)
}

func (n *Notifier) Close() {
	n.senderWaitGroup.Wait()
	defer n.fetcher.Close()
	defer func() {
		if n.dryRun {
			return
		}
		teacherService := model.NewTeacherService(n.db)
		for teacherID, lessons := range n.fetchedLessons {
			if teacher, ok := n.teachers[teacherID]; ok {
				if err := teacherService.CreateOrUpdate(teacher); err != nil {
					logger.App.Error(
						"teacherService.CreateOrUpdate failed in Notifier.Close",
						zap.Error(err), zap.Uint("teacherID", uint(teacherID)),
					)
				}
			}
			if _, err := n.lessonService.UpdateLessons(lessons); err != nil {
				logger.App.Error(
					"lessonService.UpdateLessons failed in Notifier.Close",
					zap.Error(err), zap.Uint("teacherID", uint(teacherID)),
				)
			}
		}
	}()
	defer func() {
		n.stopwatch.Stop()
		//logger.App.Info("Stopwatch report", zap.String("report", watch.Report()))
		//fmt.Println("--- stopwatch ---")
		//fmt.Println(n.stopwatch.Report())
	}()
}
