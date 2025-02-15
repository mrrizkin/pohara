package model

type JtUserRole struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
	RoleID uint `json:"role_id" gorm:"primaryKey"`
}

func (JtUserRole) TableName() string {
	return "jt_user_role"
}
