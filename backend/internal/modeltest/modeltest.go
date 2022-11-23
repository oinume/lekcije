package modeltest

import (
	"fmt"
	"time"

	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/randoms"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func NewFollowingTeacher(setters ...func(ft *model2.FollowingTeacher)) *model2.FollowingTeacher {
	followingTeacher := &model2.FollowingTeacher{}
	for _, setter := range setters {
		setter(followingTeacher)
	}
	if followingTeacher.UserID == 0 {
		followingTeacher.UserID = uint(randoms.MustNewInt64(10000000))
	}
	if followingTeacher.TeacherID == 0 {
		followingTeacher.TeacherID = uint(randoms.MustNewInt64(10000000))
	}
	return followingTeacher
}

func NewLesson(setters ...func(l *model2.Lesson)) *model2.Lesson {
	lesson := &model2.Lesson{}
	for _, setter := range setters {
		setter(lesson)
	}
	if lesson.TeacherID == 0 {
		lesson.TeacherID = uint(randoms.MustNewInt64(10000000))
	}
	if lesson.Datetime.IsZero() {
		dt := time.Now().Add(4 * time.Hour)
		lesson.Datetime = time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), 0, 0, 0, jst)
	}
	if lesson.Status == "" {
		lesson.Status = "available"
	}
	return lesson
}

func NewNotificationTimeSpan(setters ...func(nts *model2.NotificationTimeSpan)) *model2.NotificationTimeSpan {
	timeSpan := &model2.NotificationTimeSpan{}
	for _, setter := range setters {
		setter(timeSpan)
	}
	if timeSpan.UserID == 0 {
		timeSpan.UserID = uint(randoms.MustNewInt64(10000000))
	}
	if timeSpan.Number == 0 {
		timeSpan.Number = uint8(randoms.MustNewInt64(255))
	}
	if timeSpan.FromTime == "" {
		timeSpan.FromTime = ""
	}
	if timeSpan.ToTime == "" {
		timeSpan.ToTime = ""
	}
	return timeSpan
}

func NewTeacher(setters ...func(t *model2.Teacher)) *model2.Teacher {
	teacher := &model2.Teacher{}
	for _, setter := range setters {
		setter(teacher)
	}
	if teacher.ID == 0 {
		teacher.ID = uint(randoms.MustNewInt64(100000))
	}
	if teacher.Name == "" {
		teacher.Name = "teacher " + randoms.MustNewString(8)
	}
	if teacher.Gender == "" {
		teacher.Gender = "female"
	}
	if teacher.Birthday.IsZero() {
		teacher.Birthday = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	//if teacher.LastLessonAt.IsZero() {
	//	teacher.LastLessonAt = time.Now().UTC().Add(72 * time.Hour)
	//}
	return teacher
}

func NewUser(setters ...func(u *model2.User)) *model2.User {
	user := &model2.User{}
	for _, setter := range setters {
		setter(user)
	}
	if user.Name == "" {
		user.Name = "lekcije taro " + randoms.MustNewString(8)
	}
	if user.Email == "" {
		user.Email = fmt.Sprintf("lekcije-%s@example.com", randoms.MustNewString(8))
	}
	if user.PlanID == 0 {
		user.PlanID = uint8(model.DefaultMPlanID)
	}
	return user
}

func NewUserAPIToken(setters ...func(uat *model2.UserAPIToken)) *model2.UserAPIToken {
	userAPIToken := &model2.UserAPIToken{}
	for _, setter := range setters {
		setter(userAPIToken)
	}
	if userAPIToken.Token == "" {
		userAPIToken.Token = randoms.MustNewString(32)
	}
	if userAPIToken.UserID == 0 {
		userAPIToken.UserID = uint(randoms.MustNewInt64(10000000))
	}
	return userAPIToken
}

func NewUserGoogle(setters ...func(ug *model2.UserGoogle)) *model2.UserGoogle {
	userGoogle := &model2.UserGoogle{}
	for _, setter := range setters {
		setter(userGoogle)
	}
	if userGoogle.GoogleID == "" {
		userGoogle.GoogleID = randoms.MustNewString(32)
	}
	return userGoogle
}
