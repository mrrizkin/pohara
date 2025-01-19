package model

type MUserAttribute struct {
	ID       uint   `json:"id"       gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	Location string `json:"location"`
	Language string `json:"language"`
}

func (MUserAttribute) TableName() string {
	return "m_user_attribute"
}
