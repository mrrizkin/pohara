package model

type MUserAttribute struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserID   uint   `json:"user_id"`
	Location string `json:"location"`
	Language string `json:"language"`
}
