package crawler

import (
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/model"
)

type teacherIDLoader interface {
	Load() ([]uint32, error)
}

type specificTeacherIDLoader struct {
	idString string
}

func (l *specificTeacherIDLoader) Load() ([]uint32, error) {
	sids := strings.Split(l.idString, ",")
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

type followedTeacherIDLoader struct {
	db *gorm.DB
}

func (l *followedTeacherIDLoader) Load() ([]uint32, error) {
	ids, err := model.NewFollowingTeacherService(l.db).FindTeacherIDs()
	if err != nil {
		return nil, err
	}
	return ids, nil
}

type scrapingOrder int

const (
	byRating scrapingOrder = iota + 1
	byNew
)

type scrapingTeacherIDLoader struct {
	order scrapingOrder
}

func (l *scrapingTeacherIDLoader) Load() ([]uint32, error) {
	panic("implement me")
}
