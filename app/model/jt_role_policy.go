package model

type JtRolePolicy struct {
	ID       uint `json:"id"        gorm:"primaryKey"`
	RoleID   uint `json:"role_id"`
	PolicyID uint `json:"policy_id"`
}

func (JtRolePolicy) TableName() string {
	return "jt_role_policy"
}
