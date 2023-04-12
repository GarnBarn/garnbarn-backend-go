package repository

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"gorm.io/gorm"
)

type Tag interface {
	Create(tag *model.Tag) error
}

type tag struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) Tag {
	// Migrate the db
	db.AutoMigrate(&model.Tag{})

	return &tag{
		db: db,
	}
}

func (t *tag) Create(tag *model.Tag) error {
	result := t.db.Create(tag)
	return result.Error
}
