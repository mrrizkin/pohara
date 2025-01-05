package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type ByteNullable sql.NullByte

func Byte(b byte) ByteNullable {
	return ByteNullable{Valid: true, Byte: b}
}

func ByteNull() ByteNullable {
	return ByteNullable{Valid: false}
}

// Scan implements the Scanner interface.
func (n *ByteNullable) Scan(value interface{}) error {
	return (*sql.NullByte)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n ByteNullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Byte, nil
}

func (n ByteNullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Byte)
	}
	return json.Marshal(nil)
}

func (n *ByteNullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Byte)
	if err == nil {
		n.Valid = true
	}
	return err
}
