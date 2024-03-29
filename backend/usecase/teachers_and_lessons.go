package usecase

import (
	"bytes"
	"fmt"
	"time"

	"github.com/oinume/lekcije/backend/model2"
)

type teachersAndLessons struct {
	data         map[uint]*model2.TeacherLessons
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

// FilterBy filters out by NotificationTimeSpanList.
// If a lesson is within NotificationTimeSpanList, it'll be included in returned value.
func (tal *teachersAndLessons) FilterBy(list model2.NotificationTimeSpanList) *teachersAndLessons {
	if len(list) == 0 {
		return tal
	}
	ret := newTeachersAndLessons(len(tal.data))
	for teacherID, tl := range tal.data {
		lessons := make([]*model2.Lesson, 0, len(tl.Lessons))
		for _, lesson := range tl.Lessons {
			dt := lesson.Datetime
			t, _ := time.Parse("15:04", fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute()))
			if list.Within(t) {
				lessons = append(lessons, lesson)
			}
		}
		ret.data[teacherID] = model2.NewTeacherLessons(tl.Teacher, lessons)
	}
	return ret
}

func (tal *teachersAndLessons) String() string {
	b := new(bytes.Buffer)
	for _, tl := range tal.data {
		_, _ = fmt.Fprintf(b, "Teacher: %+v", tl.Teacher)
		_, _ = fmt.Fprint(b, ", Lessons:")
		for _, l := range tl.Lessons {
			_, _ = fmt.Fprintf(b, " {%+v}", l)
		}
	}
	return b.String()
}

func newTeachersAndLessons(length int) *teachersAndLessons {
	return &teachersAndLessons{
		data:         make(map[uint]*model2.TeacherLessons, length),
		lessonsCount: -1,
		teacherIDs:   make([]uint32, 0, length),
	}
}
