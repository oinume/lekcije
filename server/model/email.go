package model

type Email struct {
	raw       string
	encrypted string
}

func NewEmailFromRaw(raw string) *Email {
	return &Email{
		raw:       raw,
		encrypted: "",
	}
}

func NewEmailFromEncrypted(encrypted string) *Email {
	return &Email{
		raw:       "",
		encrypted: encrypted,
	}
}

func (e *Email) String() string {
	return e.Raw()
}

func (e *Email) Raw() string {
	return e.raw
}

func (e *Email) Encrypted() string {
	return e.encrypted
}

// Scan implements the Scanner interface.
// The value type must be time.Time or string / []byte (formatted time-string),
// otherwise Scan fails.
//func (e *Email) Scan(value interface{}) error {
//	if value == nil {
//		e.raw = ""
//		e.encrypted = ""
//		return nil
//	}
//
//	switch v := value.(type) {
//	case []byte:
//		encrypted := string(v)
//		nt.Time, err = parseDateTime(string(v), time.UTC)
//		nt.Valid = (err == nil)
//		return
//	case string:
//		nt.Time, err = parseDateTime(v, time.UTC)
//		nt.Valid = (err == nil)
//		return
//	}
//
//	nt.Valid = false
//	return fmt.Errorf("Can't convert %T to time.Time", value)
//}

// Value implements the driver Valuer interface.
//func (e *Email) Value() (driver.Value, error) {
//	if e.encrypted == "" {
//		return nil, nil
//	}
//	return e.encrypted, nil
//}

func (e *Email) encrypt() {
	//e.encrypted = ""
}

func (e *Email) decrypt() {
	//e.raw = ""
}
