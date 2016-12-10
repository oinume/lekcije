package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/errors"
)

const teacherUrlBase = "http://eikaiwa.dmm.com/teacher/index/%v/"

var (
	idsRegexp        = regexp.MustCompile(`^[\d,]+$`)
	teacherUrlRegexp = regexp.MustCompile(`https?://eikaiwa.dmm.com/teacher/index/([\d]+)`)
)

type Teacher struct {
	ID        uint32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
