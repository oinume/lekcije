package usecase

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/internal/slice_util"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/util"
)

var (
	bugUserIDs    = []uint{3725, 3698}
	bugTeacherIDs = []uint{32202, 9521, 10134, 42129} //nolint
)

type Notifier struct {
	appLogger                   *zap.Logger
	db                          *gorm.DB
	errorRecorder               *ErrorRecorder
	fetcher                     repository.LessonFetcher
	dryRun                      bool
	lessonUsecase               *Lesson
	notificationTimeSpanUsecase *NotificationTimeSpan
	statNotifierUsecase         *StatNotifier
	teacherUsecase              *Teacher
	teachers                    map[uint]*model2.Teacher
	fetchedLessons              map[uint][]*model2.Lesson
	sender                      emailer.Sender
	senderWaitGroup             *sync.WaitGroup
	followingTeacherRepo        repository.FollowingTeacher
	sync.Mutex
}

func NewNotifier(
	appLogger *zap.Logger,
	db *gorm.DB,
	errorRecorder *ErrorRecorder,
	fetcher repository.LessonFetcher,
	dryRun bool,
	lessonUsecase *Lesson,
	notificationTimeSpanUsecase *NotificationTimeSpan,
	statNotifierUsecase *StatNotifier,
	teacherUsecase *Teacher,
	sender emailer.Sender,
	followingTeacherRepo repository.FollowingTeacher,
) *Notifier {
	return &Notifier{
		appLogger:                   appLogger,
		db:                          db,
		errorRecorder:               errorRecorder,
		fetcher:                     fetcher,
		dryRun:                      dryRun,
		lessonUsecase:               lessonUsecase,
		notificationTimeSpanUsecase: notificationTimeSpanUsecase,
		statNotifierUsecase:         statNotifierUsecase,
		teacherUsecase:              teacherUsecase,
		teachers:                    make(map[uint]*model2.Teacher, 1000),
		fetchedLessons:              make(map[uint][]*model2.Lesson, 1000),
		sender:                      sender,
		senderWaitGroup:             &sync.WaitGroup{},
		followingTeacherRepo:        followingTeacherRepo,
	}
}

func (n *Notifier) SendNotification(ctx context.Context, user *model2.User) error {
	const maxFetchErrorCount = 5
	teacherIDs, err := n.followingTeacherRepo.FindTeacherIDsByUserID(
		ctx,
		user.ID,
		maxFetchErrorCount,
		time.Now().Add(-1*60*24*time.Hour), /* 2 months */
	)
	if err != nil {
		return err
	}
	if slices.Contains(bugUserIDs, user.ID) {
		n.appLogger.Info(
			"SendNotification(bug)",
			zap.Int("userID", int(user.ID)),
			zap.String("teacherIDs", strings.Join(slice_util.Map(teacherIDs, func(v uint, _ int) string {
				return fmt.Sprint(v)
			}), ",")),
		)
	}

	if len(teacherIDs) == 0 {
		return nil
	}

	// Comment out due to papertrail limit
	//logger.App.Info("n", zap.Uint("userID", uint(user.ID)), zap.Int("teachers", len(teacherIDs)))

	availableTeachersAndLessons := newTeachersAndLessons(1000)
	wg := &sync.WaitGroup{}
	for _, teacherID := range teacherIDs {
		wg.Add(1)
		go func(teacherID uint) {
			defer wg.Done()
			fetched, newAvailable, err := n.fetchAndExtractNewAvailableLessons(ctx, teacherID)
			if err != nil {
				if errors.IsNotFound(err) {
					if err := n.teacherUsecase.IncrementFetchErrorCount(ctx, teacherID, 1); err != nil {
						n.appLogger.Error(
							"IncrementFetchErrorCount failed",
							zap.Uint("teacherID", teacherID), zap.Error(err),
						)
					}
					n.appLogger.Warn("Cannot find teacher", zap.Uint("teacherID", teacherID))
				}
				// TODO: Record a case https://eikaiwa.dmm.com is down
				n.appLogger.Error("Cannot fetch teacher", zap.Uint("teacherID", teacherID), zap.Error(err))
				return
			}

			n.Lock()
			defer n.Unlock()
			n.teachers[teacherID] = fetched.Teacher
			if _, ok := n.fetchedLessons[teacherID]; !ok {
				n.fetchedLessons[teacherID] = make([]*model2.Lesson, 0, 500)
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

	timeSpans, err := n.notificationTimeSpanUsecase.FindByUserID(ctx, user.ID)
	if err != nil {
		return err
	}
	filteredAvailable := availableTeachersAndLessons.FilterBy(timeSpans)
	if err := n.sendNotificationToUser(ctx, user, filteredAvailable); err != nil {
		return err
	}

	ctx, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "Notifier.SendNotification.sleep")
	span.SetAttributes(attribute.KeyValue{
		Key:   "userID",
		Value: attribute.Int64Value(int64(user.ID)),
	})
	defer span.End()

	time.Sleep(150 * time.Millisecond)

	return nil
}

// Returns teacher, fetchedLessons, newAvailableLessons, error
func (n *Notifier) fetchAndExtractNewAvailableLessons(
	ctx context.Context,
	teacherID uint,
) (*model2.TeacherLessons, *model2.TeacherLessons, error) {
	ctx, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "NotificationTimeSpanService.FindByUserID")
	span.SetAttributes(attribute.KeyValue{
		Key:   "teacherID",
		Value: attribute.Int64Value(int64(teacherID)),
	})
	defer span.End()

	teacher, fetchedLessons, err := n.fetcher.Fetch(ctx, teacherID)
	if err != nil {
		return nil, nil, err
	}
	n.appLogger.Debug(
		"fetcher.Fetch",
		zap.Uint("teacherID", teacher.ID),
		zap.Int("lessons", len(fetchedLessons)),
	)

	//fmt.Printf("fetchedLessons ---\n")
	//for _, l := range fetchedLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}

	now := time.Now().In(config.LocalLocation())
	fromDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.LocalLocation())
	toDate := fromDate.Add(24 * 6 * time.Hour)
	lastFetchedLessons, err := n.lessonUsecase.FindLessons(ctx, teacher.ID, fromDate, toDate)
	//lastFetchedLessons, err := n.lessonService.FindLessons(ctx, uint32(teacher.ID), fromDate, toDate)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Printf("lastFetchedLessons ---\n")
	//for _, l := range lastFetchedLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}

	//modelTeacher, modelFetchedLessons := n.toModel(teacher, fetchedLessons)
	//newAvailableLessons := n.lessonService.GetNewAvailableLessons(ctx, lastFetchedLessons, modelFetchedLessons)
	newAvailableLessons := n.lessonUsecase.GetNewAvailableLessons(ctx, lastFetchedLessons, fetchedLessons)
	//fmt.Printf("newAvailableLessons ---\n")
	//for _, l := range newAvailableLessons {
	//	fmt.Printf("teacherID=%v, datetime=%v, status=%v\n", l.TeacherId, l.Datetime, l.Status)
	//}
	return model2.NewTeacherLessons(teacher, fetchedLessons),
		model2.NewTeacherLessons(teacher, newAvailableLessons),
		nil
}

func (n *Notifier) sendNotificationToUser(
	ctx context.Context,
	user *model2.User,
	lessonsByTeacher *teachersAndLessons,
) error {
	ctx, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "Notifier.sendNotificationToUser")
	defer span.End()
	if slices.Contains(bugUserIDs, user.ID) {
		n.appLogger.Info(
			"sendNotificationToUser(bug)",
			zap.Uint("userID", user.ID),
			zap.String("teacherIDs", strings.Join(slice_util.Map(lessonsByTeacher.teacherIDs, func(v uint, _ int) string {
				return fmt.Sprint(v)
			}), ",")),
			zap.Int("countLessons", lessonsByTeacher.CountLessons()),
		)
	}

	lessonsCount := 0
	var teacherIDs []int
	for teacherID, l := range lessonsByTeacher.data {
		teacherIDs = append(teacherIDs, int(teacherID))
		lessonsCount += len(l.Lessons)
	}
	if lessonsByTeacher.CountLessons() == 0 {
		// Don't send notification
		return nil
	}

	sort.Ints(teacherIDs)
	var teacherIDs2 []uint
	var teacherNames []string
	for _, id := range teacherIDs {
		teacherIDs2 = append(teacherIDs2, uint(id))
		teacherNames = append(teacherNames, n.teachers[uint(id)].Name)
	}

	// TODO: getEmailTemplate as a static file
	t := emailer.NewTemplate("notifier", getEmailTemplateJP())
	data := struct {
		To                string
		TeacherNames      string
		TeacherIDs        []uint
		Teachers          map[uint]*model2.Teacher
		LessonsPerTeacher map[uint]*model2.TeacherLessons
		WebURL            string
	}{
		To:                user.Email,
		TeacherNames:      strings.Join(teacherNames, ", "),
		TeacherIDs:        teacherIDs2,
		Teachers:          n.teachers,
		LessonsPerTeacher: lessonsByTeacher.data,
		WebURL:            config.WebURL(),
	}
	email, err := emailer.NewEmailFromTemplate(t, data)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to create emailer.Email from template: to=%v", user.Email),
		)
	}
	email.SetCustomArg("email_type", model2.EmailTypeNewLessonNotifier)
	email.SetCustomArg("user_id", fmt.Sprint(user.ID))
	email.SetCustomArg("teacher_ids", strings.Join(util.UintToStringSlice(teacherIDs2...), ","))
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
			n.errorRecorder.Record(ctx, err, fmt.Sprint(user.ID))
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
<a href="https://eikaiwa.dmm.com/teacher/index/{{ $teacherID }}/">PC</a>
<a href="https://eikaiwa.dmm.com/teacher/schedule/{{ $teacherID }}/">Mobile</a>

{{ end }}

お知らせ ─────────────────
Patreonによるサポートプログラムを開始しました。詳しくは下記をご覧ください。また、すでにサポートして下さっている皆さま、ありがとうございます。
https://lekcije.amebaownd.com/posts/10340104
─────────────────────────

空きレッスンの通知の解除は<a href="{{ .WebURL }}/me">こちら</a>

<a href="https://goo.gl/forms/CIGO3kpiQCGjtFD42">お問い合わせ</a>
	`)
}

func (n *Notifier) Close(ctx context.Context, stat *model2.StatNotifier) {
	n.senderWaitGroup.Wait()
	defer n.fetcher.Close()
	defer func() {
		if n.dryRun {
			return
		}
		ctx, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "lessonService.UpdateLessons")
		defer span.End()

		for teacherID, lessons := range n.fetchedLessons {
			if teacher, ok := n.teachers[teacherID]; ok {
				if err := n.teacherUsecase.CreateOrUpdate(ctx, teacher); err != nil {
					n.appLogger.Error(
						"teacherService.CreateOrUpdate failed in Notifier.Close",
						zap.Error(err), zap.Uint("teacherID", teacherID),
					)
					n.errorRecorder.Record(ctx, err, "")
				}
			}
			if _, err := n.lessonUsecase.UpdateLessons(ctx, lessons); err != nil {
				n.appLogger.Error(
					"lessonService.UpdateLessons failed in Notifier.Close",
					zap.Error(err), zap.Uint("teacherID", teacherID),
				)
				n.errorRecorder.Record(ctx, err, "")
			}
		}
	}()
	if stat.Interval == 10 {
		stat.Elapsed = uint(time.Now().UTC().Sub(stat.Datetime) / time.Millisecond)
		stat.FollowedTeacherCount = uint(len(n.teachers))
		if err := n.statNotifierUsecase.CreateOrUpdate(ctx, stat); err != nil {
			n.appLogger.Error("statNotifierUsecase.CreateOrUpdate failed", zap.Error(err))
		}
	}
}
