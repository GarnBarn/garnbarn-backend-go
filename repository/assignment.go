package repository

import (
	"context"
	"time"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AssignmentRepository interface {
	GetAllAssignment(author string, formPresent bool) ([]model.Assignment, error)
	CreateAssignment(assignment *model.Assignment) error
	GetByID(id int) (*model.Assignment, error)
	Update(assignment *model.Assignment) error
	DeleteAssignment(assignmentId int) error
}

type assignmentRepository struct {
	db                *gorm.DB
	repositoryContext context.Context
	tagRepository     Tag
}

func NewAssignmentRepository(db *gorm.DB, repositoryContext context.Context, tagRepository Tag) AssignmentRepository {
	// Migrate the db
	db.AutoMigrate(&model.Assignment{})

	return &assignmentRepository{
		db:                db,
		repositoryContext: repositoryContext,
		tagRepository:     tagRepository,
	}
}

func (a *assignmentRepository) CreateAssignment(assignmentData *model.Assignment) error {
	logrus.Debug("Executing Create on %T", assignmentData)

	res := a.db.WithContext(a.repositoryContext).Create(assignmentData)

	// HandleError
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		logrus.Error(res.Error)
		return res.Error
	}

	tagData, err := a.tagRepository.GetByID(assignmentData.TagID)
	if err != nil {
		return err
	}

	assignmentData.Tag = tagData
	return nil
}

func (a *assignmentRepository) GetAllAssignment(author string, fromPresent bool) (result []model.Assignment, err error) {
	now := time.Now()

	baseQuery := a.db.Model(&model.Assignment{}).Joins("Tag").WithContext(a.repositoryContext).Where("assignments.author = ?", author)

	var res *gorm.DB
	if fromPresent {
		res = baseQuery.Where("due_date >= ?", now.Unix()*1000).Find(&result)
	} else {
		res = baseQuery.Find(&result)
	}

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (a *assignmentRepository) GetByID(id int) (*model.Assignment, error) {
	var result model.Assignment
	response := a.db.WithContext(a.repositoryContext).First(&result, id)
	return &result, response.Error
}

func (a *assignmentRepository) Update(assignment *model.Assignment) error {
	result := a.db.WithContext(a.repositoryContext).Save(assignment)
	return result.Error
}

func (a *assignmentRepository) DeleteAssignment(assignmentId int) error {
	logrus.Info("Delete assignment an id: ", assignmentId)
	result := a.db.WithContext(a.repositoryContext).Delete(&model.Assignment{}, assignmentId)
	return result.Error
}
