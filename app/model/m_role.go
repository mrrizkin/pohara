package model

// Role represents a set of permissions
type MRole struct {
	Name        string `gorm:"uniqueIndex"`
	Description string
}
