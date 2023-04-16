package repository

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AssignmentRepository interface {
	CreateAssignment(assignment *model.Assignment) error
	DeleteAssignment(assignmentId int) error
}

type assignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	// Migrate the db
	db.AutoMigrate(&model.Assignment{})

	return &assignmentRepository{
		db: db,
	}
}

func (a *assignmentRepository) CreateAssignment(assignmentData *model.Assignment) error {
	logrus.Debug("Executing Create on %T", assignmentData)

	res := a.db.Create(assignmentData)

	// HandleError
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		logrus.Error(res.Error)
		return res.Error
	}

	a.db.Joins("Tag").First(assignmentData, assignmentData.ID)
	return nil
}

func (a *assignmentRepository) DeleteAssignment(assignmentId int) error {
	logrus.Info("Delete assignment an id: ", assignmentId)
	result := a.db.Joins("Tag").Delete(&model.Assignment{}, assignmentId)
	return result.Error
}
