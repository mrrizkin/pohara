package model

type JtRolePolicy struct {
	RoleID   uint `json:"role_id"   gorm:"primaryKey"`
	PolicyID uint `json:"policy_id" gorm:"primaryKey"`
}

func (JtRolePolicy) TableName() string {
	return "jt_role_policy"
}
