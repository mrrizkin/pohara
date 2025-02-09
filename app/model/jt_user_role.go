package model

type JtUserRole struct {
	ID     uint `json:"id"      gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}

func (JtUserRole) TableName() string {
	return "jt_user_role"
}
