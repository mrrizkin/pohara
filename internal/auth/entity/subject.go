package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

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
