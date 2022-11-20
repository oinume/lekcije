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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/infrastructure/dmm_eikaiwa"
	"github.com/oinume/lekcije/backend/internal/mock"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
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
	lessonUsecase := usecase.NewLesson(repos.Lesson(), repos.LessonStatusLog())

	fetcherMockTransport, err := mock.NewHTMLTransport("../infrastructure/dmm_eikaiwa/testdata/3986.html")
	if err != nil {
		t.Fatalf("fetcher.NewMockTransport failed: err=%v", err)
	}
	fetcherHTTPClient := &http.Client{
		Transport: fetcherMockTransport,
	}

	t.Run("10_users", func(t *testing.T) {
		var users []*model.User
		const numOfUsers = 10
		for i := 0; i < numOfUsers; i++ {
			name := fmt.Sprintf("oinume+%02d", i)
			user := helper.CreateUser(t, name, name+"@gmail.com")
			teacher := helper.CreateRandomTeacher(t)
			helper.CreateFollowingTeacher(t, user.ID, teacher)
			users = append(users, user)
		}

		mCountryList := registry.MustNewMCountryList(context.Background(), db.DB())
		fetcher := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)
		senderTransport := &mockSenderTransport{}
		senderHTTPClient := &http.Client{
			Transport: senderTransport,
		}
		sender := emailer.NewSendGridSender(senderHTTPClient, appLogger)
		n := usecase.NewNotifier(appLogger, db, errorRecorder, fetcher, true, lessonUsecase, sender, nil)

		ctx := context.Background()
		for _, user := range users {
			if err := n.SendNotification(ctx, user); err != nil {
				t.Fatalf("SendNotification failed: err=%v", err)
			}
		}
		// Wait all async requests are done
		n.Close(ctx, &model.StatNotifier{
			Datetime:             time.Now().UTC(),
			Interval:             10,
			Elapsed:              1000,
			UserCount:            uint32(len(users)),
			FollowedTeacherCount: uint32(len(users)),
		})

		//if got, want := senderTransport.called, numOfUsers; got <= want {
		//	t.Errorf("unexpected senderTransport.called: got=%v, want=%v", got, want)
		//}
	})

	t.Run("narrow_down_with_notification_time_span", func(t *testing.T) {
		user := helper.CreateRandomUser(t)
		teacher := helper.CreateRandomTeacher(t)
		helper.CreateFollowingTeacher(t, user.ID, teacher)

		notificationTimeSpanService := model.NewNotificationTimeSpanService(helper.DB(t))
		timeSpans := []*model.NotificationTimeSpan{
			{UserID: user.ID, Number: 1, FromTime: "02:00:00", ToTime: "03:00:00"},
			{UserID: user.ID, Number: 2, FromTime: "06:00:00", ToTime: "07:00:00"},
		}
		if err := notificationTimeSpanService.UpdateAll(user.ID, timeSpans); err != nil {
			t.Fatalf("UpdateAll failed: err=%v", err)
		}

		errorRecorder := usecase.NewErrorRecorder(appLogger, &repository.NopErrorRecorder{})
		mCountryList := registry.MustNewMCountryList(context.Background(), db.DB())
		fetcher := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)
		senderTransport := &mockSenderTransport{}
		senderHTTPClient := &http.Client{
			Transport: senderTransport,
		}
		sender := emailer.NewSendGridSender(senderHTTPClient, appLogger)
		n := usecase.NewNotifier(appLogger, db, errorRecorder, fetcher, true, lessonUsecase, sender, nil)
		if err := n.SendNotification(context.Background(), user); err != nil {
			t.Fatalf("SendNotification failed: err=%v", err)
		}

		n.Close(context.Background(), &model.StatNotifier{
			Datetime:             time.Now().UTC(),
			Interval:             10,
			Elapsed:              1000,
			UserCount:            1,
			FollowedTeacherCount: 1,
		}) // Wait all async requests are done before reading request body
		content := senderTransport.requestBody
		// TODO: table drive test
		if !strings.Contains(content, "02:30") {
			t.Errorf("content must contain 02:30 due to notification time span")
		}
		if !strings.Contains(content, "06:00") {
			t.Errorf("content must contain 06:00 due to notification time span")
		}
		if strings.Contains(content, "05:00") {
			t.Errorf("content must not contain 23:30 due to notification time span")
		}
		//fmt.Printf("content = %v\n", content)
	})
}

func TestNotifier_Close(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	db := helper.DB(t)
	appLogger := logger.NewAppLogger(os.Stdout, zapcore.DebugLevel)
	repos := mysqltest.NewRepositories(db.DB())
	lessonUsecase := usecase.NewLesson(repos.Lesson(), repos.LessonStatusLog())

	senderTransport := &mockSenderTransport{}
	senderHTTPClient := &http.Client{
		Transport: senderTransport,
	}
	sender := emailer.NewSendGridSender(senderHTTPClient, appLogger)

	user := helper.CreateRandomUser(t)
	teacher := helper.CreateTeacher(t, 3982, "Hena")
	helper.CreateFollowingTeacher(t, user.ID, teacher)

	errorRecorder := usecase.NewErrorRecorder(appLogger, &repository.NopErrorRecorder{})
	fetcherMockTransport, err := mock.NewHTMLTransport("../infrastructure/dmm_eikaiwa/testdata/3986.html")
	r.NoError(err, "fetcher.NewMockTransport failed")
	fetcherHTTPClient := &http.Client{
		Transport: fetcherMockTransport,
	}
	mCountryList := registry.MustNewMCountryList(context.Background(), db.DB())
	fetcher := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)

	n := usecase.NewNotifier(appLogger, db, errorRecorder, fetcher, false, lessonUsecase, sender, nil)
	err = n.SendNotification(context.Background(), user)
	r.NoError(err, "SendNotification failed")
	n.Close(context.Background(), &model.StatNotifier{
		Datetime:             time.Now().UTC(),
		Interval:             10,
		Elapsed:              1000,
		UserCount:            1,
		FollowedTeacherCount: 1,
	})

	teacherService := model.NewTeacherService(db)
	updatedTeacher, err := teacherService.FindByPK(teacher.ID)
	r.NoError(err)
	a.NotEqual(teacher.CountryID, updatedTeacher.CountryID)
	a.NotEqual(teacher.FavoriteCount, updatedTeacher.FavoriteCount)
	a.NotEqual(teacher.Rating, updatedTeacher.Rating)
	a.NotEqual(teacher.ReviewCount, updatedTeacher.ReviewCount)
}

func Test_Notifier_All(t *testing.T) {
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	appLogger := logger.NewAppLogger(os.Stdout, zapcore.DebugLevel)
	errorRecorder := usecase.NewErrorRecorder(appLogger, &repository.NopErrorRecorder{})
	lessonUsecase := usecase.NewLesson(repos.Lesson(), repos.LessonStatusLog())
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
	sender := emailer.NewSendGridSender(senderHTTPClient, appLogger)
	notifier1 := usecase.NewNotifier(appLogger, db, errorRecorder, fetcher1, false, lessonUsecase, sender, nil)

	user := helper.CreateRandomUser(t)
	teacher := helper.CreateTeacher(t, 49393, "Judith")
	helper.CreateFollowingTeacher(t, user.ID, teacher)

	ctx := context.Background()
	if err := notifier1.SendNotification(ctx, user); err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}
	notifier1.Close(ctx, &model.StatNotifier{
		Datetime:             time.Now().UTC(),
		Interval:             10,
		Elapsed:              1000,
		UserCount:            1,
		FollowedTeacherCount: 1,
	})

	time.Sleep(1 * time.Second)

	fetcher2 := dmm_eikaiwa.NewLessonFetcher(fetcherHTTPClient, 1, false, mCountryList, appLogger)
	notifier2 := usecase.NewNotifier(appLogger, db, errorRecorder, fetcher2, false, lessonUsecase, sender, nil)
	if err := notifier2.SendNotification(ctx, user); err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}
	notifier2.Close(ctx, &model.StatNotifier{
		Datetime:             time.Now().UTC(),
		Interval:             10,
		Elapsed:              2000,
		UserCount:            1,
		FollowedTeacherCount: 1,
	})
}

func createUsersAndFollowingTeachers(t *testing.T, num int) []*model.User {
	users := make([]*model.User, 0, num)
	for i := 0; i < num; i++ {
		name := fmt.Sprintf("oinume+%02d", i)
		user := helper.CreateUser(t, name, name+"@gmail.com")
		teacher := helper.CreateRandomTeacher(t)
		helper.CreateFollowingTeacher(t, user.ID, teacher)
		users = append(users, user)
	}
	return users
}
