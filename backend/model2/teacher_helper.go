package model2

import (
	"regexp"
	"strconv"

	"github.com/oinume/lekcije/backend/errors"
)

var (
	idRegexp         = regexp.MustCompile(`^[\d]+$`)
	teacherURLRegexp = regexp.MustCompile(`https?://eikaiwa.dmm.com/teacher/index/([\d]+)`)
)

func NewTeacher(id uint) *Teacher {
	return &Teacher{ID: id}
}

func NewTeacherFromIDOrURL(idOrURL string) (*Teacher, error) {
	if idRegexp.MatchString(idOrURL) {
		id, err := strconv.ParseUint(idOrURL, 10, 32)
		if err != nil {
			return nil, err
		}
		return NewTeacher(uint(id)), nil
	} else if group := teacherURLRegexp.FindStringSubmatch(idOrURL); len(group) > 0 {
		id, err := strconv.ParseUint(group[1], 10, 32)
		if err != nil {
			return nil, err
		}
		return NewTeacher(uint(id)), nil
	}
	return nil, errors.NewInvalidArgumentError(
		errors.WithMessage("Failed to parse idOrURL"),
		errors.WithResource(errors.NewResource("teacher", "idOrURL", idOrURL)),
	)
}
