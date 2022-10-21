package model2

import (
	"fmt"
	"time"
)

const (
	lessonTimeFormat = "2006-01-02 15:04"
)

func (l *Lesson) String() string {
	return fmt.Sprintf(
		"TeacherID=%v, Datetime=%v, Status=%v",
		l.TeacherID, l.Datetime.Format(lessonTimeFormat), l.Status,
	)
}

type LessonDatetime time.Time

func ParseLessonDatetime(s string) (LessonDatetime, error) {
	t, err := time.Parse(lessonTimeFormat, s)
	return LessonDatetime(t), err
}

func (ld LessonDatetime) String() string {
	return time.Time(ld).Format(lessonTimeFormat)
}

type TeacherLessons struct {
	Teacher *Teacher
	Lessons []*Lesson
}

func NewTeacherLessons(t *Teacher, l []*Lesson) *TeacherLessons {
	return &TeacherLessons{Teacher: t, Lessons: l}
}
