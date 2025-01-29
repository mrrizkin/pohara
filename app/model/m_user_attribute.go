package model

type MUserAttribute struct {
	ID       uint   `json:"id"       db:"id"`
	UserID   uint   `json:"user_id"  db:"user_id"`
	Location string `json:"location" db:"location"`
	Language string `json:"language" db:"language"`
}

func (MUserAttribute) TableName() string {
	return "m_user_attribute"
}
