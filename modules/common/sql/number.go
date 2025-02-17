package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type (
	Float64Nullable sql.NullFloat64
	Int16Nullable   sql.NullInt16
	Int32Nullable   sql.NullInt32
	Int64Nullable   sql.NullInt64
)

func Float64(f float64) Float64Nullable {
	return Float64Nullable{Valid: true, Float64: f}
}

func Float64Null() Float64Nullable {
	return Float64Nullable{Valid: false}
}

func Int16(i int16) Int16Nullable {
	return Int16Nullable{Valid: true, Int16: i}
}

func Int16Null() Int16Nullable {
	return Int16Nullable{Valid: false}
}

func Int32(i int32) Int32Nullable {
	return Int32Nullable{Valid: true, Int32: i}
}

func Int32Null() Int32Nullable {
	return Int32Nullable{Valid: false}
}

func Int64(i int64) Int64Nullable {
	return Int64Nullable{Valid: true, Int64: i}
}

func Int64Null() Int64Nullable {
	return Int64Nullable{Valid: false}
}

// Scan implements the Scanner interface.
func (n *Float64Nullable) Scan(value interface{}) error {
	return (*sql.NullFloat64)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n Float64Nullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

func (n Float64Nullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}
	return json.Marshal(nil)
}

func (n *Float64Nullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Float64)
	if err == nil {
		n.Valid = true
	}
	return err
}

// Scan implements the Scanner interface.
func (n *Int16Nullable) Scan(value interface{}) error {
	return (*sql.NullInt16)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n Int16Nullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int16, nil
}

func (n Int16Nullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int16)
	}
	return json.Marshal(nil)
}

func (n *Int16Nullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Int16)
	if err == nil {
		n.Valid = true
	}
	return err
}

// Scan implements the Scanner interface.
func (n *Int32Nullable) Scan(value interface{}) error {
	return (*sql.NullInt32)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n Int32Nullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int32, nil
}

func (n Int32Nullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int32)
	}
	return json.Marshal(nil)
}

func (n *Int32Nullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Int32)
	if err == nil {
		n.Valid = true
	}
	return err
}

// Scan implements the Scanner interface.
func (n *Int64Nullable) Scan(value interface{}) error {
	return (*sql.NullInt64)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n Int64Nullable) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

func (n Int64Nullable) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

func (n *Int64Nullable) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Int64)
	if err == nil {
		n.Valid = true
	}
	return err
}
