package model

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type MCountry struct {
	ID        uint16 `gorm:"primary_key"`
	Name      string
	NameJA    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*MCountry) TableName() string {
	return "m_country"
}

type MCountryService struct {
	db           *gorm.DB
	allCountries map[string]*MCountry
	sync.Mutex
}

func NewMCountryService(db *gorm.DB) *MCountryService {
	return &MCountryService{db: db}
}

func (s *MCountryService) LoadAll() (*MCountries, error) {
	values := make([]*MCountry, 0, 1000)
	sql := `SELECT * FROM m_country ORDER BY name`
	if err := s.db.Raw(sql).Scan(&values).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to load from m_country"),
		)
	}
	return NewMCountries(values), nil
}

func (s *MCountryService) GetMCountryByNameJA(nameJA string) (*MCountry, error) {
	if len(s.allCountries) == 0 {
		return nil, fmt.Errorf("m_country not loaded")
	}
	if c, ok := s.allCountries[nameJA]; ok {
		return c, nil
	} else {
		return nil, errors.NewNotFoundError(errors.WithMessagef("No MCountries for %v", nameJA))
	}
}

func NewMCountries(values []*MCountry) *MCountries {
	c := &MCountries{
		byNameJA: make(map[string]*MCountry, 1000),
	}
	for _, v := range values {
		c.byNameJA[v.NameJA] = v
	}
	return c
}

type MCountries struct {
	byNameJA map[string]*MCountry
}

func (mc *MCountries) GetByNameJA(name string) (*MCountry, bool) {
	c, ok := mc.byNameJA[name]
	return c, ok
}
