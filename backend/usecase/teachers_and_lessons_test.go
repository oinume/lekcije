package usecase

import (
	"testing"
	"time"

	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

func Test_TeachersAndLessons_FilterBy(t *testing.T) {
	user := modeltest.NewUser()
	timeSpans := []*model.NotificationTimeSpan{
		{UserID: uint32(user.ID), Number: 1, FromTime: "15:30:00", ToTime: "16:30:00"},
		{UserID: uint32(user.ID), Number: 2, FromTime: "20:00:00", ToTime: "22:00:00"},
	}
	teacher := modeltest.NewTeacher()
	// TODO: table driven test
	lessons := []*model2.Lesson{
		{TeacherID: teacher.ID, Datetime: time.Date(2018, 1, 1, 15, 0, 0, 0, time.UTC)}, // excluded
		{TeacherID: teacher.ID, Datetime: time.Date(2018, 1, 1, 16, 0, 0, 0, time.UTC)}, // included
		{TeacherID: teacher.ID, Datetime: time.Date(2018, 1, 1, 17, 0, 0, 0, time.UTC)}, // excluded
		{TeacherID: teacher.ID, Datetime: time.Date(2018, 1, 1, 21, 0, 0, 0, time.UTC)}, // included
		{TeacherID: teacher.ID, Datetime: time.Date(2018, 1, 1, 23, 0, 0, 0, time.UTC)}, // excluded
	}
	tal := newTeachersAndLessons(10)
	tal.data[teacher.ID] = &model2.TeacherLessons{Teacher: teacher, Lessons: lessons}

	filtered := tal.FilterBy(timeSpans)
	if got, want := filtered.CountLessons(), 2; got != want {
		t.Fatalf("unexpected filtered lessons count: got=%v, want=%v", got, want)
	}

	wantTimes := []struct {
		hour, minute int
	}{
		{16, 0},
		{21, 0},
	}
	tl := filtered.data[teacher.ID]
	for i, wantTime := range wantTimes {
		if got, want := tl.Lessons[i].Datetime.Hour(), wantTime.hour; got != want {
			t.Errorf("unexpected hour: got=%v, want=%v", got, want)
		}
		if got, want := tl.Lessons[i].Datetime.Minute(), wantTime.minute; got != want {
			t.Errorf("unexpected minute: got=%v, want=%v", got, want)
		}
	}
}

func Test_TeachersAndLessons_FilterByEmpty(t *testing.T) {
	timeSpans := make([]*model.NotificationTimeSpan, 0)
	teacher := modeltest.NewTeacher()
	// TODO: table driven test
	lessons := []*model2.Lesson{
		{TeacherID: teacher.ID, Datetime: time.Date(2018, 1, 1, 15, 0, 0, 0, time.UTC)},
		{TeacherID: teacher.ID, Datetime: time.Date(2018, 1, 1, 16, 0, 0, 0, time.UTC)},
	}
	tal := newTeachersAndLessons(10)
	tal.data[teacher.ID] = &model2.TeacherLessons{Teacher: teacher, Lessons: lessons}

	filtered := tal.FilterBy(timeSpans)
	if got, want := filtered.CountLessons(), len(lessons); got != want {
		t.Fatalf("unexpected filtered lessons count: got=%v, want=%v", got, want)
	}

	wantTimes := []struct {
		hour, minute int
	}{
		{15, 0},
		{16, 0},
	}
	tl := filtered.data[teacher.ID]
	for i, wantTime := range wantTimes {
		if got, want := tl.Lessons[i].Datetime.Hour(), wantTime.hour; got != want {
			t.Errorf("unexpected hour: got=%v, want=%v", got, want)
		}
		if got, want := tl.Lessons[i].Datetime.Minute(), wantTime.minute; got != want {
			t.Errorf("unexpected minute: got=%v, want=%v", got, want)
		}
	}
}
