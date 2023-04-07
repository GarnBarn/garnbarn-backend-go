package model

import "gorm.io/gorm"

type ExampleDB struct {
	gorm.Model
	Name  string
	Value int
}
