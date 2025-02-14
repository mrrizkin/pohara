package access

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// Action represents possible actions on resources
type Action struct {
	Valid    bool
	Type     string
	Resource string
}

func NewResourceAction(resource string, action string) Action {
	return Action{
		Valid:    true,
		Type:     action,
		Resource: resource,
	}
}

func NewViewPageAction(resource string) Action {
	return Action{
		Valid:    true,
		Type:     "view",
		Resource: resource,
	}
}

// Scan implements the Scanner interface.
func (n *Action) Scan(value interface{}) error {
	if value == nil {
		n.Type, n.Resource, n.Valid = "", "", false
		return nil
	}

	err := convertAssign(n, value)
	n.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (n Action) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String(), nil
}

func (n Action) String() string {
	if !n.Valid {
		return ""
	}
	return n.Type + "::" + n.Resource
}

func (n Action) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String())
	}

	return json.Marshal(nil)
}

func (n *Action) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Type, n.Resource, n.Valid = "", "", false
		return nil
	}

	var action string
	if err := json.Unmarshal(b, &action); err != nil {
		n.Valid = false
		return err
	}
	err := convertAssign(n, action)
	n.Valid = err == nil
	return err
}

func convertAssign(a *Action, value any) error {
	var s string

	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("invalid type for Action: %T", value)
	}

	// Split the string by "::" separator
	parts := strings.Split(s, "::")
	if len(parts) != 2 {
		return fmt.Errorf("invalid action format: %s", s)
	}

	// Assign the parts to the Action struct
	a.Type = parts[0]
	a.Resource = parts[1]

	return nil
}
