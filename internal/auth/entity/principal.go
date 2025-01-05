package entity

import (
	"encoding/json"

	"github.com/mrrizkin/pohara/internal/common/nanoid"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint             `json:"id"          gorm:"primary_key"`
	CreatedAt   sql.TimeNullable `json:"created_at"`
	UpdatedAt   sql.TimeNullable `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at"  gorm:"index"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
}

func (Role) TableName() string {
	return "m_roles"
}

type Subject interface {
	GetSubjectID() uint
	GetSubjectName() string
	GetAttributes() map[string]interface{}
}

type BaseSubject struct {
	gorm.Model
	PublicID string `json:"public_id"`
	Name     string `json:"name"`
}

func (u *BaseSubject) BeforeCreate(tx *gorm.DB) (err error) {
	u.PublicID = nanoid.New()
	return
}

func (s BaseSubject) GetSubjectID() uint {
	return s.ID
}

func (s BaseSubject) GetSubjectName() string {
	return s.Name
}

func (s BaseSubject) GetAttributes() map[string]interface{} {
	attrs := make(map[string]interface{})
	bytes, _ := json.Marshal(s)
	json.Unmarshal(bytes, &attrs)
	return attrs
}
