package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type BoolNullable sql.NullBool

func Bool(b bool) BoolNullable {
	return BoolNullable{Valid: true, Bool: b}
}

func BoolNull() BoolNullable {
	return BoolNullable{Valid: false}
}

// Scan implements the Scanner interface.
func (n *BoolNullable) Scan(value interface{}) error {
	return (*sql.NullBool)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n BoolNullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

func (n BoolNullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

func (n *BoolNullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Bool)
	if err == nil {
		n.Valid = true
	}
	return err
}
