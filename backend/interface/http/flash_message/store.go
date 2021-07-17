package flash_message

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
)

type Store interface {
	Load(key string) (*FlashMessage, error)
	Save(value *FlashMessage) error
}

type StoreMySQL struct {
	db *gorm.DB
}

func NewStoreMySQL(db *gorm.DB) *StoreMySQL {
	return &StoreMySQL{db: db}
}

func (s *StoreMySQL) Load(key string) (*FlashMessage, error) {
	flashMessage := &model.FlashMessage{}
	if result := s.db.First(flashMessage, model.FlashMessage{ID: key}); result.Error != nil {
		if result.RecordNotFound() {
			return nil, errors.NewAnnotatedError(errors.CodeNotFound)
		}
		return nil, result.Error
	}
	if flashMessage.IsExpired(time.Now().UTC()) {
		return nil, fmt.Errorf("FlashMessage is expired: key=%v", key)
	}
	v := &FlashMessage{}
	if err := json.Unmarshal([]byte(flashMessage.Value), v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *StoreMySQL) Save(value *FlashMessage) error {
	return s.SaveWithExpiration(value, time.Hour*24)
}

func (s *StoreMySQL) SaveWithExpiration(value *FlashMessage, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	fm := &model.FlashMessage{
		ID:        value.Key,
		Value:     string(bytes),
		ExpiredAt: time.Now().UTC().Add(expiration),
	}
	if result := s.db.Create(fm); result.Error != nil {
		return result.Error
	}
	return nil
}
