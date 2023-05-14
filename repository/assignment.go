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

	// Get new assignment
	newAssignment, err := a.GetByID(int(assignmentData.ID))
	if err != nil {
		logrus.Error(err)
		return err
	}
	*assignmentData = *newAssignment
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

	for _, item := range result {
		if item.TagID == 0 || item.Tag == nil {
			continue
		}

		decryptedTagData, err := a.tagRepository.GetByID(item.TagID)
		if err != nil {
			return result, err
		}
		*item.Tag = *decryptedTagData
	}

	return result, nil
}

func (a *assignmentRepository) GetByID(id int) (*model.Assignment, error) {
	var result model.Assignment
	response := a.db.WithContext(a.repositoryContext).First(&result, id)

	tagData, err := a.tagRepository.GetByID(result.TagID)
	if err != nil {
		return nil, err
	}

	result.Tag = tagData

	return &result, response.Error
}

func (a *assignmentRepository) Update(assignment *model.Assignment) error {
	result := a.db.WithContext(a.repositoryContext).Save(assignment)
	if result.Error != nil {
		return result.Error
	}

	newAssignment, err := a.GetByID(int(assignment.ID))
	if err != nil {
		return err
	}

	*assignment = *newAssignment

	return nil
}

func (a *assignmentRepository) DeleteAssignment(assignmentId int) error {
	logrus.Info("Delete assignment an id: ", assignmentId)
	result := a.db.WithContext(a.repositoryContext).Delete(&model.Assignment{}, assignmentId)
	return result.Error
}
