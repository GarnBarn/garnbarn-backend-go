package repository

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Example interface {
	Example(exampleData *model.ExampleDB) error
}

type example struct {
	db *gorm.DB
}

func NewExampleRepository(db *gorm.DB) Example {
	// Migrate the db
	db.AutoMigrate(&model.ExampleDB{})

	return &example{
		db: db,
	}
}

func (e *example) Example(exampleData *model.ExampleDB) error {
	result := e.db.Create(exampleData)
	logrus.Debug(result.Error)
	logrus.Debug(exampleData.ID)

	return result.Error
}
