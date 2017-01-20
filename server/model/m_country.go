package model

import "time"

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
