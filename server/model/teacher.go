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

const teacherUrlBase = "http://eikaiwa.dmm.com/teacher/index/%v/"

var (
	idsRegexp        = regexp.MustCompile(`^[\d,]+$`)
	teacherUrlRegexp = regexp.MustCompile(`https?://eikaiwa.dmm.com/teacher/index/([\d]+)`)
)

type Teacher struct {
	ID                uint32
	Name              string
	CountryID         uint16
	Gender            string
	Birthday          time.Time
	YearsOfExperience uint8
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

func NewTeachersFromIDsOrURL(idsOrUrl string) ([]*Teacher, error) {
	if idsRegexp.MatchString(idsOrUrl) {
		ids := strings.Split(idsOrUrl, ",")
		teachers := make([]*Teacher, 0, len(ids))
		for _, sid := range ids {
			if id, err := strconv.ParseUint(sid, 10, 32); err == nil {
				teachers = append(teachers, NewTeacher(uint32(id)))
			} else {
				continue
			}
		}
		return teachers, nil
	} else if group := teacherUrlRegexp.FindStringSubmatch(idsOrUrl); len(group) > 0 {
		id, _ := strconv.ParseUint(group[1], 10, 32)
		return []*Teacher{NewTeacher(uint32(id))}, nil
	}
	return nil, errors.InvalidArgumentf("Failed to parse idsOrUrl: %s", idsOrUrl)
}

func (t *Teacher) URL() string {
	return fmt.Sprintf(teacherUrlBase, t.ID)
}

type TeacherService struct {
	db *gorm.DB
}

func NewTeacherService(db *gorm.DB) *TeacherService {
	return &TeacherService{db: db}
}

func (s *TeacherService) CreateOrUpdate(t *Teacher) error {
	sql := fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", t.TableName())
	sql += " ON DUPLICATE KEY UPDATE"
	sql += " country_id=?, gender=?, years_of_experience=?, birthday=?"
	now := time.Now()
	values := []interface{}{
		t.ID,
		t.Name,
		t.CountryID,
		t.Gender,
		t.Birthday.Format("2006-01-02"),
		t.YearsOfExperience,
		t.FetchErrorCount,
		now.Format("2006-01-02 15:04:05"),
		now.Format("2006-01-02 15:04:05"),
		// update
		t.CountryID,
		t.Gender,
		t.YearsOfExperience,
		t.Birthday.Format("2006-01-02"),
	}

	if err := s.db.Exec(sql, values...).Error; err != nil {
		return errors.InternalWrapf(err, "Failed to INSERT or UPDATE teacher: id=%v", t.ID)
	}
	return nil
}

func (s *TeacherService) FindByPK(id uint32) (*Teacher, error) {
	teacher := &Teacher{}
	if err := s.db.First(teacher, &Teacher{ID: id}).Error; err != nil {
		return nil, errors.InternalWrapf(err, "Failed to FindByPK")
	}
	return teacher, nil
}

func (s *TeacherService) IncrementFetchErrorCount(id uint32, value int) error {
	sql := fmt.Sprintf(
		`UPDATE %s SET fetch_error_count = fetch_error_count + ?, updatd_at = NOW() WHERE id = ?`,
		new(Teacher).TableName(),
	)
	if err := s.db.Exec(sql, value, id).Error; err != nil {
		return errors.InternalWrapf(err, "Failed to UPDATE teacher: id=%v", id)
	}
	return nil
}

func (s *TeacherService) ResetFetchErrorCount(id uint32) error {
	sql := fmt.Sprintf(
		`UPDATE %s SET fetch_error_count = 0, updated_at = NOW() WHERE id = ?`,
		new(Teacher).TableName(),
	)
	if err := s.db.Exec(sql, id).Error; err != nil {
		return errors.InternalWrapf(err, "Failed to UPDATE teacher: id=%v", id)
	}
	return nil
}

func (s *TeacherService) FindByFetchErrorCountGt(count uint32) ([]*Teacher, error) {
	values := make([]*Teacher, 0, 3000)
	sql := fmt.Sprintf(`SELECT * FROM teacher WHERE fetch_error_count > ? ORDER BY id LIMIT 3000`)
	if result := s.db.Raw(sql, count).Scan(&values); result.Error != nil {
		return nil, errors.InternalWrapf(result.Error, "Failed to FindByFetchErrorCountGt")
	}
	return values, nil
}
