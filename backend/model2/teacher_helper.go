package model2

import (
	"regexp"
	"strconv"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/v4/types"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
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

func NewTeacherFromModel(t *model.Teacher) *Teacher {
	rating := decimal.New(0, 2)
	rating.SetFloat64(float64(t.Rating))
	return &Teacher{
		ID:                uint(t.ID),
		Name:              t.Name,
		CountryID:         int16(t.CountryID),
		Gender:            t.Gender,
		Birthday:          t.Birthday,
		YearsOfExperience: int8(t.YearsOfExperience),
		FavoriteCount:     uint(t.FavoriteCount),
		ReviewCount:       uint(t.ReviewCount),
		Rating: types.NullDecimal{
			Big: rating,
		},
		LastLessonAt:    t.LastLessonAt,
		FetchErrorCount: t.FetchErrorCount,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}
