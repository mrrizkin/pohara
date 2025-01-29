package model

// Role represents a set of permissions
type MRole struct {
	ID          uint   `json:"id"          db:"id"`
	Name        string `json:"name"        db:"name"`
	Description string `json:"description" db:"description"`
}

func (MRole) TableName() string {
	return "m_role"
}
