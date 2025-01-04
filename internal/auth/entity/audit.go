package entity

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog represents an authorization decision log
type AuditAuthLog struct {
	gorm.Model
	TraceID         string
	Timestamp       time.Time
	SubjectID       uint
	SubjectName     string
	Action          Action
	ResourceType    string
	ResourceID      uint
	Effect          bool
	PolicyIDs       []uint `gorm:"type:json"`
	MatchedPolicyID *uint
	Context         string `gorm:"type:json"`
	Decision        string
	Duration        int64
}

func (AuditAuthLog) TableName() string {
	return "audit_auth_logs"
}
