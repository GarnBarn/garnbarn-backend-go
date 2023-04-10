package model

import "gorm.io/gorm"

type Assignment struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	DueDate     int    `json:"dueDate"`
}
