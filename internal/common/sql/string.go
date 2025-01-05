package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type StringNullable sql.NullString

func String(s string) StringNullable {
	return StringNullable{Valid: true, String: s}
}

func StringNull() StringNullable {
	return StringNullable{Valid: false}
}

// Scan implements the Scanner interface.
func (n *StringNullable) Scan(value interface{}) error {
	return (*sql.NullString)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n StringNullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

func (n StringNullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n *StringNullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.String)
	if err == nil {
		n.Valid = true
	}
	return err
}
