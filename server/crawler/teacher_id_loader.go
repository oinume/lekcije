package crawler

import (
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/model"
)

type TeacherIDLoader interface {
	Load() ([]uint32, error)
}

type SpecificTeacherIDLoader struct {
	IDString string
}

func (l *SpecificTeacherIDLoader) Load() ([]uint32, error) {
	sids := strings.Split(l.IDString, ",")
	ids := make([]uint32, 0, len(sids))
	for _, id := range sids {
		i, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint32(i))
	}
	return ids, nil
}

type FollowedTeacherIDLoader struct {
	DB *gorm.DB
}

func (l *FollowedTeacherIDLoader) Load() ([]uint32, error) {
	ids, err := model.NewFollowingTeacherService(l.DB).FindTeacherIDs()
	if err != nil {
		return nil, err
	}
	return ids, nil
}

type ScrapingTeacherIDLoader struct {
}

func (l *ScrapingTeacherIDLoader) Load() ([]uint32, error) {
	panic("implement me")
}
