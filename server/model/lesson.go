package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/goenum"
	"github.com/oinume/lekcije/server/errors"
)

const (
	lessonTimeFormat = "2006-01-02 15:04"
)

type Lesson struct {
	TeacherID uint32    `gorm:"primary_key"`
	Datetime  time.Time `gorm:"primary_key"`
	Status    string    // TODO: enum
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Lesson) TableName() string {
	return "lesson"
}

func (l *Lesson) String() string {
	return fmt.Sprintf(
		"TeacherID=%v, Datetime=%v, Status=%v",
		l.TeacherID, l.Datetime.Format(lessonTimeFormat), l.Status,
	)
}

type LessonStatus struct {
	Finished  int `goenum:"終了"`
	Reserved  int `goenum:"予約済"`
	Available int `goenum:"予約可"`
	Cancelled int `goenum:"休講"`
}

var LessonStatuses = goenum.EnumerateStruct(&LessonStatus{
	Finished:  1,
	Reserved:  2,
	Available: 3,
	Cancelled: 4,
})

type LessonService struct {
	db *gorm.DB
}

func NewLessonService(db *gorm.DB) *LessonService {
	return &LessonService{db: db}
}

func (s *LessonService) TableName() string {
	return (&Lesson{}).TableName()
}

func (s *LessonService) UpdateLessons(lessons []*Lesson) (int64, error) {
	if len(lessons) == 0 {
		return 0, nil
	}

	updatedAt := time.Now().UTC()
	sql := fmt.Sprintf("INSERT INTO %s VALUES", s.TableName())
	values := []interface{}{}
	for _, l := range lessons {
		sql += " (?, ?, ?, ?, ?),"
		values = append(
			values,
			l.TeacherID, l.Datetime.Format(dbDatetimeFormat), strings.ToLower(l.Status), // TODO: enum?
			updatedAt.Format(dbDatetimeFormat), updatedAt.Format(dbDatetimeFormat),
		)
	}
	sql = strings.TrimSuffix(sql, ",")
	sql += " ON DUPLICATE KEY UPDATE status=VALUES(status), updated_at=VALUES(updated_at)"

	result := s.db.Exec(sql, values...)
	if err := result.Error; err != nil {
		return 0, errors.InternalWrapf(err, "")
	}

	return result.RowsAffected, nil
}

func (s *LessonService) FindLessons(teacherID uint32, fromDate, toDate time.Time) ([]*Lesson, error) {
	lessons := make([]*Lesson, 0, 1000)
	sql := strings.TrimSpace(fmt.Sprintf(`
SELECT * FROM %s
WHERE
  teacher_id = ?
  AND DATE(datetime) BETWEEN ? AND ?
ORDER BY datetime ASC
LIMIT 1000
	`, s.TableName()))

	toDateAdded := toDate.Add(24 * 2 * time.Hour)
	result := s.db.Raw(sql, teacherID, fromDate.Format("2006-01-02"), toDateAdded.Format("2006-01-02")).Scan(&lessons)
	if result.Error != nil {
		if result.RecordNotFound() {
			return lessons, nil
		}
		return nil, errors.InternalWrapf(result.Error, "Failed to find lessons: teacherID=%v", teacherID)
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to find lessons: teacherID=%v"),
			errors.WithResources(errors.NewResource(s.TableName(), "teacherID", fmt.Sprint(teacherID))),
		)
	}

	return lessons, nil
}

func (s *LessonService) GetNewAvailableLessons(oldLessons, newLessons []*Lesson) []*Lesson {
	// Pattern
	// 2016-01-01 00:00@Any -> Available
	oldLessonsMap := make(map[string]*Lesson, len(oldLessons))
	newLessonsMap := make(map[string]*Lesson, len(newLessons))
	availableLessons := make([]*Lesson, 0, len(oldLessons)+len(newLessons))
	availableLessonsMap := make(map[string]*Lesson, len(oldLessons)+len(newLessons))
	for _, l := range oldLessons {
		oldLessonsMap[l.Datetime.Format(lessonTimeFormat)] = l
	}
	for _, l := range newLessons {
		newLessonsMap[l.Datetime.Format(lessonTimeFormat)] = l
	}
	for datetime, oldLesson := range oldLessonsMap {
		if newLesson, ok := newLessonsMap[datetime]; ok && strings.ToLower(oldLesson.Status) != "available" && strings.ToLower(newLesson.Status) == "available" {
			// exists in oldLessons and newLessons
			availableLessons = append(availableLessons, newLesson)
			availableLessonsMap[datetime] = newLesson
		}
	}

	for _, l := range newLessons {
		datetime := l.Datetime.Format(lessonTimeFormat)
		if _, ok := oldLessonsMap[datetime]; !ok && strings.ToLower(l.Status) == "available" {
			// not exists in oldLessons
			availableLessons = append(availableLessons, l)
			availableLessonsMap[datetime] = l
		}
	}

	// TODO: sort availableLessonsMap by datetime
	return availableLessons
}
