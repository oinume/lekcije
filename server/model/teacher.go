package model

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const teacherUrlBase = "http://eikaiwa.dmm.com/teacher/index/%v/"

var (
	idRegexp         = regexp.MustCompile(`^[\d]+$`)
	teacherUrlRegexp = regexp.MustCompile(`https?://eikaiwa.dmm.com/teacher/index/([\d]+)`)
)

type Teacher struct {
	Id        uint32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Teacher) TableName() string {
	return "teacher"
}

func NewTeacher(id uint32) *Teacher {
	return &Teacher{Id: id}
}

func NewTeacherFromIdOrUrl(idOrUrl string) (*Teacher, error) {
	if idRegexp.MatchString(idOrUrl) {
		id, _ := strconv.ParseUint(idOrUrl, 10, 32)
		return NewTeacher(uint32(id)), nil
	} else if group := teacherUrlRegexp.FindStringSubmatch(idOrUrl); len(group) > 0 {
		id, _ := strconv.ParseUint(group[1], 10, 32)
		return NewTeacher(uint32(id)), nil
	}
	return nil, fmt.Errorf("Failed to parse idOrUrl: %s", idOrUrl)
}

func (t *Teacher) Url() string {
	return fmt.Sprintf(teacherUrlBase, t.Id)
}
