package model

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	ID            uint   `gorm:"primarykey" json:"id"`
	Name          string `json:"name" validate:"required"`
	Author        string `json:"author"`
	Color         string `json:"color,omitempty"`
	ReminderTime  string `json:"reminderTime"`
	Subscriber    string `json:"subscribe"`
	SecretKeyTotp string `json:"secretKeyTotp"`
}
