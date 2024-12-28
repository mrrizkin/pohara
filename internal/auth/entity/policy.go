package entity

import (
	"encoding/json"

	"github.com/mrrizkin/pohara/internal/common/sql"
	"gorm.io/gorm"
)

type Policies []Policy

type Policy struct {
	ID           uint             `json:"id"            gorm:"primary_key"`
	CreatedAt    sql.TimeNullable `json:"created_at"`
	UpdatedAt    sql.TimeNullable `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"deleted_at"    gorm:"index"`
	RoleID       uint             `json:"role_id"`
	ResourceType string           `json:"resource_type"` // e.g., "document", "user"
	Action       string           `json:"action"`        // e.g., "create", "read", "update", "delete"
	Effect       string           `json:"effect"`        // "allow" or "deny"
	Conditions   string           `json:"conditions"`    // JSON string of conditions
	Department   string           `json:"department"`    // "*" for all departments
}

type OptimizePolicy struct {
	ID           uint             `json:"id"            gorm:"primary_key"`
	CreatedAt    sql.TimeNullable `json:"created_at"`
	UpdatedAt    sql.TimeNullable `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"deleted_at"    gorm:"index"`
	RoleID       uint             `json:"role_id"`
	ResourceType string           `json:"resource_type"` // e.g., "document", "user"
	Action       string           `json:"action"`        // e.g., "create", "read", "update", "delete"
	Effect       string           `json:"effect"`        // "allow" or "deny"
	Conditions   []Condition      `json:"conditions"`    // JSON string of conditions
	Department   string           `json:"department"`    // "*" for all departments
}

type Condition struct {
	Type      sql.StringNullable `json:"type"`
	Attribute sql.StringNullable `json:"attribute"`
	Operator  sql.StringNullable `json:"operator"`
	Value     sql.StringNullable `json:"value"`
}

func (p Policies) Optimize() ([]*OptimizePolicy, error) {
	optimizePolicies := make([]*OptimizePolicy, len(p))
	for i, policy := range p {
		optimizePolicy, err := policy.Optimize()
		if err != nil {
			continue
		}

		optimizePolicies[i] = optimizePolicy
	}

	return optimizePolicies, nil
}

func (p *Policy) Optimize() (*OptimizePolicy, error) {
	optimizePolicy := OptimizePolicy{
		ID:           p.ID,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		DeletedAt:    p.DeletedAt,
		RoleID:       p.RoleID,
		ResourceType: p.ResourceType,
		Action:       p.Action,
		Effect:       p.Effect,
		Department:   p.Department,
	}

	var conditions []Condition
	if err := json.Unmarshal([]byte(p.Conditions), &conditions); err != nil {
		return nil, err
	}

	optimizePolicy.Conditions = conditions

	return &optimizePolicy, nil
}
