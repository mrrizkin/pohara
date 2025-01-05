package repository

import (
	"github.com/mrrizkin/pohara/internal/auth/access"
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"github.com/mrrizkin/pohara/internal/ports"
	"gorm.io/gorm"
)

type Repository struct {
	db ports.Database
}

type RepositoryDependencies struct {
	Database ports.Database
}

func NewRepository(deps RepositoryDependencies) *Repository {
	return &Repository{
		db: deps.Database,
	}
}

func (r *Repository) CreateAuditAuthLog(auditAuthLog *entity.AuditAuthLog) error {
	return r.db.Create(auditAuthLog)
}

func (r *Repository) GetRolePolicy(subject entity.Subject, resource entity.Resource, action access.Action) ([]entity.Policy, error) {
	var policies []entity.Policy
	db := r.db.GetDB().(*gorm.DB)
	err := db.Preload("Roles").
		Joins("JOIN role_policies ON role_policies.policy_id = policies.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_policies.role_id").
		Where("user_roles.subject_id = ? AND policies.resource_type = ? AND policies.action = ?",
			subject.GetSubjectID(), resource.GetResourceType(), action).
		Order("policies.priority DESC").
		Find(&policies).Error
	return policies, err
}

func (r *Repository) SearchAuditLogs(query entity.AuditQuery) (*ports.FindResult, error) {
	where := sql.Where()

	if query.StartTime != nil {
		where.And("timestamp >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		where.And("timestamp <= ?", query.EndTime)
	}
	if query.SubjectID != nil {
		where.And("subject_id = ?", *query.SubjectID)
	}
	if query.ResourceType != nil {
		where.And("resource_type = ?", *query.ResourceType)
	}
	if query.Action != nil {
		where.And("action = ?", *query.Action)
	}
	if query.Effect != nil {
		where.And("effect = ?", *query.Effect)
	}
	if query.PolicyID != nil {
		where.And("matched_policy_id = ?", *query.PolicyID)
	}

	var logs []entity.AuditAuthLog
	result, err := r.db.Find(&logs, ports.Pagination{
		Limit:  sql.Int64(int64(query.PageSize)),
		Offset: sql.Int64(int64((query.Page - 1) * query.PageSize)),
		Sort:   sql.String("timestamp DESC"),
	}, where.Cond()...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
