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

type Resource interface {
	GetResourceType() string
	GetResourceID() uint
	GetAttributes() map[string]interface{}
}

// BaseResource provides default implementation for resources
type BaseResource struct {
	gorm.Model
	ResourceType string
}

func (r BaseResource) GetResourceType() string {
	return r.ResourceType
}

func (r BaseResource) GetResourceID() uint {
	return r.ID
}

func (r BaseResource) GetAttributes() map[string]interface{} {
	attrs := make(map[string]interface{})
	bytes, _ := json.Marshal(r)
	json.Unmarshal(bytes, &attrs)
	return attrs
}
