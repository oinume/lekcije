package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

const teacherUrlBase = "https://eikaiwa.dmm.com/teacher/index/%v/"

var (
	defaultLastLessonAt = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	idsRegexp           = regexp.MustCompile(`^[\d,]+$`)
	teacherURLRegexp    = regexp.MustCompile(`https?://eikaiwa.dmm.com/teacher/index/([\d]+)`)
)

type Teacher struct {
	ID                uint32
	Name              string
	CountryID         uint16
	Gender            string
	Birthday          time.Time
	YearsOfExperience uint8
	FavoriteCount     uint32
	ReviewCount       uint32
	Rating            float32
	LastLessonAt      time.Time
	FetchErrorCount   uint8
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (*Teacher) TableName() string {
	return "teacher"
}

func NewTeacher(id uint32) *Teacher {
	return &Teacher{ID: id}
}

func NewTeachersFromIDsOrURL(idsOrURL string) ([]*Teacher, error) {
	if idsRegexp.MatchString(idsOrURL) {
		ids := strings.Split(idsOrURL, ",")
		teachers := make([]*Teacher, 0, len(ids))
		for _, sid := range ids {
			if id, err := strconv.ParseUint(sid, 10, 32); err == nil {
				teachers = append(teachers, NewTeacher(uint32(id)))
			} else {
				continue
			}
		}
		return teachers, nil
	} else if group := teacherURLRegexp.FindStringSubmatch(idsOrURL); len(group) > 0 {
		id, _ := strconv.ParseUint(group[1], 10, 32)
		return []*Teacher{NewTeacher(uint32(id))}, nil
	}
	return nil, errors.NewInvalidArgumentError(
		errors.WithMessage("Failed to parse idsOrURL"),
		errors.WithResource(errors.NewResource("teacher", "idsOrURL", idsOrURL)),
	)
}

func (t *Teacher) URL() string {
	return fmt.Sprintf(teacherUrlBase, t.ID)
}

func (t *Teacher) IsJapanese() bool {
	return t.CountryID == 392
}

type TeacherLessons struct {
	Teacher *Teacher
	Lessons []*Lesson
}

func NewTeacherLessons(t *Teacher, l []*Lesson) *TeacherLessons {
	return &TeacherLessons{Teacher: t, Lessons: l}
}

type TeacherService struct {
	db *gorm.DB
}

func NewTeacherService(db *gorm.DB) *TeacherService {
	return &TeacherService{db: db}
}

func (s *TeacherService) CreateOrUpdate(t *Teacher) error {
	sql := fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", t.TableName())
	sql += " ON DUPLICATE KEY UPDATE"
	sql += " country_id=?, gender=?, years_of_experience=?, birthday=?, favorite_count=?, review_count=?, rating=?"

	if t.LastLessonAt.IsZero() {
		t.LastLessonAt = defaultLastLessonAt
	}
	now := time.Now().UTC()
	values := []interface{}{
		t.ID,
		t.Name,
		t.CountryID,
		t.Gender,
		t.Birthday.Format(dbDateFormat),
		t.YearsOfExperience,
		t.FavoriteCount,
		t.ReviewCount,
		t.Rating,
		t.LastLessonAt.Format(dbDatetimeFormat),
		t.FetchErrorCount,
		now.Format(dbDatetimeFormat),
		now.Format(dbDatetimeFormat),
		// update
		t.CountryID,
		t.Gender,
		t.YearsOfExperience,
		t.Birthday.Format(dbDateFormat),
		t.FavoriteCount,
		t.ReviewCount,
		t.Rating,
	}

	if !t.LastLessonAt.Equal(defaultLastLessonAt) {
		// UPDATE last_lesson_at only when it's NOT default value
		sql += ", last_lesson_at=?"
		values = append(values, t.LastLessonAt.Format(dbDatetimeFormat))
	}

	if err := s.db.Exec(sql, values...).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to INSERT or UPDATE teacher"),
			errors.WithResource(errors.NewResource("teacher", "id", t.ID)),
		)
	}
	return nil
}

func (s *TeacherService) FindByPK(id uint32) (*Teacher, error) {
	teacher := &Teacher{}
	if err := s.db.First(teacher, &Teacher{ID: id}).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to FindByPK"),
			errors.WithResource(errors.NewResource(teacher.TableName(), "id", id)),
		)
	}
	return teacher, nil
}

func (s *TeacherService) IncrementFetchErrorCount(id uint32, value int) error {
	tableName := new(Teacher).TableName()
	sql := fmt.Sprintf(
		`UPDATE %s SET fetch_error_count = fetch_error_count + ?, updated_at = NOW() WHERE id = ?`,
		tableName,
	)
	if err := s.db.Exec(sql, value, id).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to UPDATE teacher"),
			errors.WithResource(errors.NewResource(tableName, "id", id)),
		)
	}
	return nil
}

func (s *TeacherService) ResetFetchErrorCount(id uint32) error {
	tableName := new(Teacher).TableName()
	sql := fmt.Sprintf(
		`UPDATE %s SET fetch_error_count = 0, updated_at = NOW() WHERE id = ?`,
		tableName,
	)
	if err := s.db.Exec(sql, id).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to UPDATE teacher"),
			errors.WithResource(errors.NewResource(tableName, "id", id)),
		)
	}
	return nil
}

func (s *TeacherService) FindByFetchErrorCountGt(count uint32) ([]*Teacher, error) {
	values := make([]*Teacher, 0, 3000)
	sql := fmt.Sprintf(`SELECT * FROM teacher WHERE fetch_error_count > ? ORDER BY id LIMIT 3000`)
	if result := s.db.Raw(sql, count).Scan(&values); result.Error != nil {
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to FindByFetchErrorCountGt"),
		)
	}
	return values, nil
}
