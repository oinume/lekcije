package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/goenum"
	"github.com/oinume/lekcije/server/errors"
)

type Lesson struct {
	TeacherId uint32    `gorm:"primary_key"`
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
		"TeacherId: %v, Datetime: %v, Status: %v",
		l.TeacherId, l.Datetime, l.Status,
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

type LessonServiceType struct {
	db *gorm.DB
}

var LessonService LessonServiceType

func (s *LessonServiceType) TableName() string {
	return (&Lesson{}).TableName()
}

func (s *LessonServiceType) UpdateLessons(lessons []*Lesson) (int64, error) {
	if len(lessons) == 0 {
		return 0, nil
	}

	now := time.Now()
	sql := fmt.Sprintf("INSERT INTO %s VALUES", s.TableName())
	values := []interface{}{}
	for _, l := range lessons {
		sql += " (?, ?, ?, ?, ?),"
		values = append(
			values,
			l.TeacherId, l.Datetime, strings.ToLower(l.Status), // TODO: enum?
			now.Format(dbDatetimeFormat), now.Format(dbDatetimeFormat),
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

func (s *LessonServiceType) FindLessons(teacherId uint32, fromDate, toDate time.Time) ([]*Lesson, error) {
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
	result := s.db.Raw(sql, teacherId, fromDate.Format("2006-01-02"), toDateAdded.Format("2006-01-02")).Scan(&lessons)
	if result.Error != nil {
		if result.RecordNotFound() {
			return lessons, nil
		}
		return nil, errors.InternalWrapf(result.Error, "Failed to find lessons: teacherId=%v", teacherId)
	}

	return lessons, nil
}

func (s *LessonServiceType) GetNewAvailableLessons(oldLessons, newLessons []*Lesson) []*Lesson {
	// Pattern
	// 2016-01-01 00:00@Any -> Available
	oldLessonsMap := make(map[time.Time]string, len(oldLessons))
	newLessonsMap := make(map[time.Time]string, len(newLessons))
	for _, l := range oldLessons {
		oldLessonsMap[l.Datetime] = strings.ToLower(l.Status)
	}
	for _, l := range newLessons {
		newLessonsMap[l.Datetime] = strings.ToLower(l.Status)
	}
	for datetime, status := range oldLessonsMap {
		if newStatus, ok := newLessonsMap[datetime]; ok {
			fmt.Printf("oldStatus = %v, newStatus = %v\n", status, newStatus)
		}
	}
	return nil
	/*
	   @classmethod
	   def get_new_reservable_schedules(
	       cls, old_schedules: List["Schedule"], new_schedules: List["Schedule"]
	   ) -> List["Schedule"]:
	       old = [o.to_json() for o in old_schedules]
	       new = [o.to_json() for o in new_schedules]

	       differ = difflib.Differ()
	       ret = []
	       diffs = list(differ.compare(old, new))
	       for i, d in enumerate(diffs):
	           if d.startswith("+ "):
	               if (i == len(diffs) - 1) or (i < len(diffs) - 1 and not diffs[i+1].startswith("? ")):
	                   schedule = cls.from_json(d[1:])
	                   if schedule.status == ScheduleStatus.reservable:
	                       ret.append(schedule)
	           #print("line{}:{}".format(i, d))
	       return ret
	*/
}
