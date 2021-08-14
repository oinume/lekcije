package notifier

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/jinzhu/gorm"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/fetcher"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/stopwatch"
	"github.com/oinume/lekcije/backend/util"
)

type Notifier struct {
	appLogger       *zap.Logger
	db              *gorm.DB
	fetcher         *fetcher.LessonFetcher
	dryRun          bool
	lessonService   *model.LessonService
	teachers        map[uint32]*model.Teacher
	fetchedLessons  map[uint32][]*model.Lesson
	sender          emailer.Sender
	senderWaitGroup *sync.WaitGroup
	stopwatch       stopwatch.Stopwatch
	storageClient   *storage.Client
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

func NewNotifier(
	appLogger *zap.Logger,
	db *gorm.DB,
	fetcher *fetcher.LessonFetcher,
	dryRun bool,
	sender emailer.Sender,
	sw stopwatch.Stopwatch,
	storageClient *storage.Client,
) *Notifier {
	if sw == nil {
		sw = stopwatch.NewSync()
	}
	return &Notifier{
		appLogger:       appLogger,
		db:              db,
		fetcher:         fetcher,
		dryRun:          dryRun,
		teachers:        make(map[uint32]*model.Teacher, 1000),
		fetchedLessons:  make(map[uint32][]*model.Lesson, 1000),
		sender:          sender,
		senderWaitGroup: &sync.WaitGroup{},
		stopwatch:       sw,
		storageClient:   storageClient,
	}
}

func (n *Notifier) SendNotification(ctx context.Context, user *model.User) error {
	followingTeacherService := model.NewFollowingTeacherService(n.db)
	n.lessonService = model.NewLessonService(n.db)
	const maxFetchErrorCount = 5
	teacherIDs, err := followingTeacherService.FindTeacherIDsByUserID(
		ctx,
		user.ID,
		maxFetchErrorCount,
		time.Now().Add(-1*60*24*time.Hour), /* 2 months */
	)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to FindTeacherIDsByUserID(): userID=%v", user.ID),
		)
	}

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
			defer wg.Done()
			fetched, newAvailable, err := n.fetchAndExtractNewAvailableLessons(ctx, teacherID)
			if err != nil {
				if errors.IsNotFound(err) {
					if err := model.NewTeacherService(n.db).IncrementFetchErrorCount(teacherID, 1); err != nil {
						n.appLogger.Error(
							"IncrementFetchErrorCount failed",
							zap.Uint("teacherID", uint(teacherID)), zap.Error(err),
						)
					}
					n.appLogger.Warn("Cannot find teacher", zap.Uint("teacherID", uint(teacherID)))
				}
				// TODO: Record a case eikaiwa.dmm.com is down
				n.appLogger.Error("Cannot fetch teacher", zap.Uint("teacherID", uint(teacherID)), zap.Error(err))
				return
			}

			n.Lock()
			defer n.Unlock()
			n.teachers[teacherID] = fetched.Teacher
			if _, ok := n.fetchedLessons[teacherID]; !ok {
				n.fetchedLessons[teacherID] = make([]*model.Lesson, 0, 500)
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
	timeSpans, err := notificationTimeSpanService.FindByUserID(ctx, user.ID)
	if err != nil {
		return err
	}
	filteredAvailable := availableTeachersAndLessons.FilterBy(model.NotificationTimeSpanList(timeSpans))
	if err := n.sendNotificationToUser(ctx, user, filteredAvailable); err != nil {
		return err
	}

	_, span := trace.StartSpan(ctx, "Notifier.SendNotification.sleep")
	defer span.End()
	span.Annotatef([]trace.Attribute{
		trace.Int64Attribute("userID", int64(user.ID)),
	}, "userID:%d", user.ID)
	time.Sleep(150 * time.Millisecond)

	return nil
}

// Returns teacher, fetchedLessons, newAvailableLessons, error
func (n *Notifier) fetchAndExtractNewAvailableLessons(
	ctx context.Context,
	teacherID uint32,
) (*model.TeacherLessons, *model.TeacherLessons, error) {
	_, span := trace.StartSpan(ctx, "Notifier.fetchAndExtractNewAvailableLessons")
	defer span.End()
	span.Annotatef([]trace.Attribute{
		trace.Int64Attribute("teacherID", int64(teacherID)),
	}, "teacherID:%d", teacherID)

	teacher, fetchedLessons, err := n.fetcher.Fetch(ctx, teacherID)
	if err != nil {
		return nil, nil, err
	}
	n.appLogger.Debug(
		"fetcher.Fetch",
		zap.Uint("teacherID", uint(teacher.ID)),
		zap.Int("lessons", len(fetchedLessons)),
	)

	//fmt.Printf("fetchedLessons ---\n")
	//for _, l := range fetchedLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}

	now := time.Now().In(config.LocalLocation())
	fromDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.LocalLocation())
	toDate := fromDate.Add(24 * 6 * time.Hour)
	lastFetchedLessons, err := n.lessonService.FindLessons(ctx, teacher.ID, fromDate, toDate)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Printf("lastFetchedLessons ---\n")
	//for _, l := range lastFetchedLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}

	newAvailableLessons := n.lessonService.GetNewAvailableLessons(ctx, lastFetchedLessons, fetchedLessons)
	//fmt.Printf("newAvailableLessons ---\n")
	//for _, l := range newAvailableLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}
	return model.NewTeacherLessons(teacher, fetchedLessons),
		model.NewTeacherLessons(teacher, newAvailableLessons),
		nil
}

func (n *Notifier) sendNotificationToUser(
	ctx context.Context,
	user *model.User,
	lessonsPerTeacher *teachersAndLessons,
) error {
	_, span := trace.StartSpan(ctx, "Notifier.sendNotificationToUser")
	defer span.End()

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

	n.appLogger.Info("sendNotificationToUser", zap.String("email", user.Email))

	n.senderWaitGroup.Add(1)
	go func(email *emailer.Email) {
		defer n.senderWaitGroup.Done()
		if err := n.sender.Send(ctx, email); err != nil {
			n.appLogger.Error(
				"Failed to sendNotificationToUser",
				zap.String("email", user.Email), zap.Error(err),
			)
			util.SendErrorToRollbar(err, fmt.Sprint(user.ID))
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

お知らせ ─────────────────
Patreonによるサポートプログラムを開始しました。詳しくは下記をご覧ください。また、すでにサポートして下さっている皆さま、ありがとうございます。
https://lekcije.amebaownd.com/posts/10340104
─────────────────────────

空きレッスンの通知の解除は<a href="{{ .WebURL }}/me">こちら</a>

<a href="https://goo.gl/forms/CIGO3kpiQCGjtFD42">お問い合わせ</a>
	`)
}

func (n *Notifier) Close(ctx context.Context, stat *model.StatNotifier) {
	n.senderWaitGroup.Wait()
	defer n.fetcher.Close()
	defer func() {
		if n.dryRun {
			return
		}
		_, span := trace.StartSpan(ctx, "lessonService.UpdateLessons")
		defer span.End()

		teacherService := model.NewTeacherService(n.db)
		for teacherID, lessons := range n.fetchedLessons {
			if teacher, ok := n.teachers[teacherID]; ok {
				if err := teacherService.CreateOrUpdate(teacher); err != nil {
					n.appLogger.Error(
						"teacherService.CreateOrUpdate failed in Notifier.Close",
						zap.Error(err), zap.Uint("teacherID", uint(teacherID)),
					)
					util.SendErrorToRollbar(err, "")
				}
			}
			if _, err := n.lessonService.UpdateLessons(lessons); err != nil {
				n.appLogger.Error(
					"lessonService.UpdateLessons failed in Notifier.Close",
					zap.Error(err), zap.Uint("teacherID", uint(teacherID)),
				)
				util.SendErrorToRollbar(err, "")
			}
		}
	}()
	// TODO: remove
	defer func() {
		n.stopwatch.Stop()
		if n.storageClient != nil {
			//fmt.Println("--- stopwatch ---")
			//fmt.Println(n.stopwatch.Report())
			if err := n.uploadStopwatchReport(); err != nil {
				n.appLogger.Error("uploadStopwatchReport failed", zap.Error(err))
			}
		}
	}()
	if stat.Interval == 10 {
		stat.Elapsed = uint32(time.Now().UTC().Sub(stat.Datetime) / time.Millisecond)
		stat.FollowedTeacherCount = uint32(len(n.teachers))
		if err := model.NewStatNotifierService(n.db).CreateOrUpdate(stat); err != nil {
			n.appLogger.Error("statNotifierService.CreateOrUpdate failed", zap.Error(err))
		}
	}
}

func (n *Notifier) uploadStopwatchReport() error {
	if n.storageClient == nil {
		return nil
	}

	path := time.Now().UTC().Format("stopwatch/20060102/150405.txt")
	w := n.storageClient.Bucket("lekcije").Object(path).NewWriter(context.Background())
	if _, err := w.Write([]byte(n.stopwatch.Report())); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}
