package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type TimeNullable sql.NullTime

func Time(t time.Time) TimeNullable {
	return TimeNullable{Valid: true, Time: t}
}

func TimeNull() TimeNullable {
	return TimeNullable{Valid: false}
}

// Scan implements the Scanner interface.
func (n *TimeNullable) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n TimeNullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n TimeNullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *TimeNullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Time)
	if err == nil {
		n.Valid = true
	}
	return err
}
