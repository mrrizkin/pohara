package entity

import (
	"time"

	"github.com/mrrizkin/pohara/internal/auth/access"
	"gorm.io/gorm"
)

// AuditLog represents an authorization decision log
type AuditAuthLog struct {
	gorm.Model
	TraceID         string
	Timestamp       time.Time
	SubjectID       uint
	SubjectName     string
	Action          access.Action
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

type AuditQuery struct {
	StartTime    *time.Time
	EndTime      *time.Time
	SubjectID    *uint
	ResourceType *string
	Action       *access.Action
	Effect       *bool
	PolicyID     *uint
	Page         int
	PageSize     int
}
