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
