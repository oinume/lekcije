package model

import (
	"database/sql/driver"
	"fmt"
	"os"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/util"
)

var encryptionKey = os.Getenv("ENCRYPTION_KEY")

type Email struct {
	raw       string
	encrypted string
}

func NewEmailFromRaw(raw string) (Email, error) {
	e := Email{
		raw: raw,
	}
	if encrypted, err := util.EncryptString(raw, encryptionKey); err == nil {
		e.encrypted = encrypted
	} else {
		return Email{}, err
	}
	return e, nil
}

func NewEmailFromEncrypted(encrypted string) (Email, error) {
	e := Email{
		encrypted: encrypted,
	}
	if decrypted, err := util.DecryptString(encrypted, encryptionKey); err == nil {
		e.raw = decrypted
	} else {
		return Email{}, err
	}
	return e, nil
}

func (e Email) String() string {
	return e.Raw()
}

func (e Email) Raw() string {
	return e.raw
}

func (e Email) Encrypted() string {
	return e.encrypted
}

// Scan implements the Scanner interface.
// The value type must be time.Time or string / []byte (formatted time-string),
// otherwise Scan fails.
func (e *Email) Scan(value interface{}) error {
	if value == nil {
		e.raw = ""
		e.encrypted = ""
		return nil
	}

	switch value.(type) {
	case []byte:
		encrypted := fmt.Sprintf("%s", value)
		if decrypted, err := util.DecryptString(encrypted, encryptionKey); err == nil {
			e.raw = decrypted
			return nil
		} else {
			return err
		}
	case string:
		if decrypted, err := util.DecryptString(value.(string), encryptionKey); err == nil {
			e.raw = decrypted
			return nil
		} else {
			return err
		}
	}
	return errors.Internalf("Cannot convert %T to model.Email", value)
}

// Value implements the driver Valuer interface.
func (e Email) Value() (driver.Value, error) {
	if e.encrypted != "" {
		return e.encrypted, nil
	}
	if e.raw != "" {
		return util.EncryptString(e.raw, encryptionKey)
	}
	return nil, nil
}
