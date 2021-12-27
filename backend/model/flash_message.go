package model

import "time"

type FlashMessage struct {
	ID        string
	Value     string
	ExpiredAt time.Time
}

func (*FlashMessage) TableName() string {
	return "flash_message"
}

func (fm *FlashMessage) IsExpired(now time.Time) bool {
	if fm.ExpiredAt.IsZero() {
		return false
	}
	return fm.ExpiredAt.Before(now)
}

// TODO: delete
