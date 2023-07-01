package usecase_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/v4/types"
	"go.uber.org/zap/zapcore"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/infrastructure/dmm_eikaiwa"
	"github.com/oinume/lekcije/backend/infrastructure/send_grid"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/mock"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/registry"
	"github.com/oinume/lekcije/backend/usecase"
)

type mockSenderTransport struct {
	sync.Mutex
	called      int
	requestBody string
}

func (t *mockSenderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	t.called++
	defer t.Unlock()
	time.Sleep(time.Millisecond * 500)
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusAccepted,
		Status:     "202 Accepted",
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return resp, err
	}
	t.requestBody = string(body)
	defer req.Body.Close()
	resp.Body = io.NopCloser(strings.NewReader(""))
	return resp, nil
}

func Test_Notifier_SendNotification(t *testing.T) {
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	appLogger := logger.NewAppLogger(os.Stdout, zapcore.DebugLevel)
	errorRecorder := usecase.NewErrorRecorder(appLogger, &repository.NopErrorRecorder{})
	notificationTimeSpanUsecase := usecase.NewNotificationTimeSpan(repos.NotificationTimeSpan())
	lessonUsecase := usecase.NewLesson(repos.Lesson(), repos.LessonStatusLog())
	statNotifierUsecase := usecase.NewStatNotifier(repos.StatNotifier())
	teacherUsecase := usecase.NewTeacher(repos.Teacher())

	fetcherMockTransport, err := mock.NewHTMLTransport("../infrastructure/dmm_eikaiwa/testdata/49393.html")
	if err != nil {
		t.Fatalf("fetcher.NewMockTransport failed: err=%v", err)
	}
	fetcherHTTPClient := &http.Client{
		Transport: fetcherMockTransport,
	}

	t.Run("10_users", func(t *testing.T) {
		ctx := context.Background()
		var users []*model2.User
		const numOfUsers = 10
		for i := 0; i < numOfUsers; i++ {
			user := modeltest.NewUser(func(u *model2.User) {
				name := fmt.Sprintf("oinume+%02d", i)
				u.Name = name
				u.Email = name + "@gmail.com"
			})
			repos.CreateUsers(ctx, t, user)
			teacher := modeltest.NewTeacher()
			repos.CreateTeachers(ctx, t, teacher)
			followingTeacher := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
				ft.UserID = user.ID
				ft.TeacherID = teacher.ID
			})
			repos.CreateFollowingTeachers(ctx, t, followingTeacher)
			users = append(users, user)
		}

		mCountryList := registry.MustNewMCountryList(context.Background(), db.DB())
		fetcher := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)
		senderTransport := &mockSenderTransport{}
		senderHTTPClient := &http.Client{
			Transport: senderTransport,
		}
		sender := send_grid.NewSendGridEmailSender(senderHTTPClient, appLogger)
		n := usecase.NewNotifier(
			appLogger, db, errorRecorder, fetcher, true, lessonUsecase, notificationTimeSpanUsecase,
			statNotifierUsecase, teacherUsecase, sender, repos.FollowingTeacher(),
		)

		for _, user := range users {
			if err := n.SendNotification(ctx, user); err != nil {
				t.Fatalf("SendNotification failed: err=%v", err)
			}
		}
		// Wait all async requests are done
		n.Close(ctx, &model2.StatNotifier{
			Datetime:             time.Now().UTC(),
			Interval:             1,
			Elapsed:              1000,
			UserCount:            uint(len(users)),
			FollowedTeacherCount: uint(len(users)),
		})

		//if got, want := senderTransport.called, numOfUsers; got <= want {
		//	t.Errorf("unexpected senderTransport.called: got=%v, want=%v", got, want)
		//}
	})

	t.Run("narrow_down_with_notification_time_span", func(t *testing.T) {
		ctx := context.Background()
		user := modeltest.NewUser()
		repos.CreateUsers(ctx, t, user)
		teacher := modeltest.NewTeacher(func(t *model2.Teacher) {
			t.ID = 49393
			t.Name = "Judith"
		})
		repos.CreateTeachers(ctx, t, teacher)
		followingTeacher := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
			ft.UserID = user.ID
			ft.TeacherID = teacher.ID
		})
		repos.CreateFollowingTeachers(ctx, t, followingTeacher)

		notificationTimeSpanService := model.NewNotificationTimeSpanService(helper.DB(t))
		timeSpans := []*model.NotificationTimeSpan{
			{UserID: uint32(user.ID), Number: 1, FromTime: "11:00:00", ToTime: "12:00:00"},
			{UserID: uint32(user.ID), Number: 2, FromTime: "16:00:00", ToTime: "17:00:00"},
		}
		if err := notificationTimeSpanService.UpdateAll(uint32(user.ID), timeSpans); err != nil {
			t.Fatalf("UpdateAll failed: err=%v", err)
		}

		mCountryList := registry.MustNewMCountryList(context.Background(), db.DB())
		fetcher := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)
		senderTransport := &mockSenderTransport{}
		senderHTTPClient := &http.Client{
			Transport: senderTransport,
		}
		sender := send_grid.NewSendGridEmailSender(senderHTTPClient, appLogger)
		n := usecase.NewNotifier(
			appLogger, db, errorRecorder, fetcher, true, lessonUsecase, notificationTimeSpanUsecase,
			statNotifierUsecase, teacherUsecase, sender, repos.FollowingTeacher(),
		)
		if err := n.SendNotification(context.Background(), user); err != nil {
			t.Fatalf("SendNotification failed: err=%v", err)
		}

		n.Close(context.Background(), &model2.StatNotifier{
			Datetime:             time.Now().UTC(),
			Interval:             2,
			Elapsed:              1000,
			UserCount:            1,
			FollowedTeacherCount: 1,
		}) // Wait all async requests are done before reading request body
		content := senderTransport.requestBody
		// TODO: table drive test
		if !strings.Contains(content, "11:30") {
			t.Errorf("content must contain 11:30 due to notification time span")
		}
		if !strings.Contains(content, "16:00") {
			t.Errorf("content must contain 16:00 due to notification time span")
		}
		if strings.Contains(content, "14:00") {
			t.Errorf("content must not contain 14:00 due to notification time span")
		}
		//fmt.Printf("content = %v\n", content)
	})
}

func TestNotifier_Close(t *testing.T) {
	ctx := context.Background()
	db := helper.DB(t)
	helper.TruncateAllTables(t)
	appLogger := logger.NewAppLogger(os.Stdout, zapcore.DebugLevel)
	repos := mysqltest.NewRepositories(db.DB())
	lessonUsecase := usecase.NewLesson(repos.Lesson(), repos.LessonStatusLog())
	notificationTimeSpanUsecase := usecase.NewNotificationTimeSpan(repos.NotificationTimeSpan())
	statNotifierUsecase := usecase.NewStatNotifier(repos.StatNotifier())
	teacherUsecase := usecase.NewTeacher(repos.Teacher())

	senderTransport := &mockSenderTransport{}
	senderHTTPClient := &http.Client{
		Transport: senderTransport,
	}
	sender := send_grid.NewSendGridEmailSender(senderHTTPClient, appLogger)

	user := modeltest.NewUser()
	repos.CreateUsers(ctx, t, user)
	teacher := modeltest.NewTeacher(func(t *model2.Teacher) {
		t.ID = 49393
		t.Name = "Judith"
		t.CountryID = int16(608)
		t.Birthday = time.Time{}
		t.YearsOfExperience = 2
		t.FavoriteCount = 559
		t.ReviewCount = 1267
		t.Rating = types.NullDecimal{Big: decimal.New(int64(498), 2)}
	})
	repos.CreateTeachers(ctx, t, teacher)
	followingTeacher := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
		ft.UserID = user.ID
		ft.TeacherID = teacher.ID
	})
	repos.CreateFollowingTeachers(ctx, t, followingTeacher)

	errorRecorder := usecase.NewErrorRecorder(appLogger, &repository.NopErrorRecorder{})
	fetcherMockTransport, err := mock.NewHTMLTransport("../infrastructure/dmm_eikaiwa/testdata/49393.html")
	if err != nil {
		t.Fatalf("fetcher.NewMockTransport failed: %v", err)
	}
	fetcherHTTPClient := &http.Client{
		Transport: fetcherMockTransport,
	}
	mCountryList := registry.MustNewMCountryList(context.Background(), db.DB())
	fetcher := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)

	n := usecase.NewNotifier(
		appLogger, db, errorRecorder, fetcher, false, lessonUsecase, notificationTimeSpanUsecase,
		statNotifierUsecase, teacherUsecase, sender, repos.FollowingTeacher(),
	)
	if err := n.SendNotification(context.Background(), user); err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}
	n.Close(context.Background(), &model2.StatNotifier{
		Datetime:             time.Now().UTC(),
		Interval:             3,
		Elapsed:              1000,
		UserCount:            1,
		FollowedTeacherCount: 1,
	})

	updatedTeacher, err := repos.Teacher().FindByID(ctx, teacher.ID)
	if err != nil {
		t.Fatalf("FindByPK failed: %v", err)
	}
	assertion.AssertEqual(t, teacher.CountryID, updatedTeacher.CountryID, "CountryID")
	assertion.AssertEqual(t, teacher.FavoriteCount, updatedTeacher.FavoriteCount, "FavoriteCount")
	assertion.AssertEqual(t, teacher.Rating.String(), updatedTeacher.Rating.String(), "Rating")
	assertion.AssertEqual(t, teacher.ReviewCount, updatedTeacher.ReviewCount, "ReviewCount")
}

func Test_Notifier_All(t *testing.T) {
	ctx := context.Background()
	db := helper.DB(t)
	helper.TruncateAllTables(t)
	repos := mysqltest.NewRepositories(db.DB())
	appLogger := logger.NewAppLogger(os.Stdout, zapcore.DebugLevel)
	errorRecorder := usecase.NewErrorRecorder(appLogger, &repository.NopErrorRecorder{})
	lessonUsecase := usecase.NewLesson(repos.Lesson(), repos.LessonStatusLog())
	notificationTimeSpanUsecase := usecase.NewNotificationTimeSpan(repos.NotificationTimeSpan())
	statNotifierUsecase := usecase.NewStatNotifier(repos.StatNotifier())
	teacherUsecase := usecase.NewTeacher(repos.Teacher())
	mCountryList := registry.MustNewMCountryList(context.Background(), db.DB())

	fetcherMockTransport := mock.NewResponseTransport(func(rt *mock.ResponseTransport, req *http.Request) *http.Response {
		resp := &http.Response{
			Header:     make(http.Header),
			Request:    req,
			StatusCode: http.StatusOK,
			Status:     "200 OK",
		}
		resp.Header.Set("Content-Type", "text/html; charset=UTF-8")

		var file string
		if rt.NumCalled == 1 {
			file = "../infrastructure/dmm_eikaiwa/testdata/49393.html"
		} else {
			file = "../infrastructure/dmm_eikaiwa/testdata/49393-reserved.html"
		}
		f, err := os.Open(file)
		if err != nil {
			t.Fatalf("Failed to open file: %v: %v", file, err)
		}
		resp.Body = f
		return resp
	})
	//fetcherMockTransport, err := mock.NewHTMLTransport("../infrastructure/dmm_eikaiwa/testdata/49393.html")
	//if err != nil {
	//	t.Fatalf("fetcher.NewMockTransport failed: err=%v", err)
	//}
	fetcherHTTPClient := &http.Client{
		Transport: fetcherMockTransport,
	}

	fetcher1 := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)
	senderTransport := &mockSenderTransport{}
	senderHTTPClient := &http.Client{
		Transport: senderTransport,
	}
	sender := send_grid.NewSendGridEmailSender(senderHTTPClient, appLogger)
	notifier1 := usecase.NewNotifier(
		appLogger, db, errorRecorder, fetcher1, false, lessonUsecase, notificationTimeSpanUsecase,
		statNotifierUsecase, teacherUsecase, sender, repos.FollowingTeacher(),
	)

	user := modeltest.NewUser()
	repos.CreateUsers(ctx, t, user)
	teacher := modeltest.NewTeacher(func(t *model2.Teacher) {
		t.ID = 49393
		t.Name = "Judith"
	})
	repos.CreateTeachers(ctx, t, teacher)
	followingTeacher := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
		ft.UserID = user.ID
		ft.TeacherID = teacher.ID
	})
	repos.CreateFollowingTeachers(ctx, t, followingTeacher)

	if err := notifier1.SendNotification(ctx, user); err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}
	notifier1.Close(ctx, &model2.StatNotifier{
		Datetime:             time.Now().UTC(),
		Interval:             4,
		Elapsed:              1000,
		UserCount:            1,
		FollowedTeacherCount: 1,
	})

	time.Sleep(1 * time.Second)

	fetcher2 := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)
	notifier2 := usecase.NewNotifier(
		appLogger, db, errorRecorder, fetcher2, false, lessonUsecase, notificationTimeSpanUsecase,
		statNotifierUsecase, teacherUsecase, sender, repos.FollowingTeacher(),
	)
	if err := notifier2.SendNotification(ctx, user); err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}
	notifier2.Close(ctx, &model2.StatNotifier{
		Datetime:             time.Now().UTC(),
		Interval:             5,
		Elapsed:              2000,
		UserCount:            1,
		FollowedTeacherCount: 1,
	})
	// TODO: check lesson is updated
}
