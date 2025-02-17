package model

import "time"

type MigrationHistory struct {
	ID        string    `json:"id"         gorm:"primary_key"`
	Batch     int       `json:"batch"`
	CreatedAt time.Time `json:"created_at"`
}

func (MigrationHistory) TableName() string {
	return "__migration_history__"
}
