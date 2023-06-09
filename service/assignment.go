package service

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/sirupsen/logrus"
)

type AssignmentService interface {
	CreateAssignment(assignment *model.Assignment) error
	GetAllAssignment(author string, fromPresent bool) ([]model.Assignment, error)
	GetAssignmentById(assignmentId int) (model.AssignmentPublic, error)
	UpdateAssignment(updateAssignmentRequest *model.UpdateAssignmentRequest, id int) (*model.Assignment, error)
	DeleteAssignment(assignmentId int) error
}

type assignmentService struct {
	assignmentRepository repository.AssignmentRepository
}

func NewAssignmentService(assignmentRepository repository.AssignmentRepository) AssignmentService {
	return &assignmentService{
		assignmentRepository: assignmentRepository,
	}
}

func (a *assignmentService) CreateAssignment(assignmentData *model.Assignment) error {
	return a.assignmentRepository.CreateAssignment(assignmentData)
}

func (a *assignmentService) GetAllAssignment(author string, fromPresent bool) ([]model.Assignment, error) {
	return a.assignmentRepository.GetAllAssignment(author, fromPresent)
}

func (a *assignmentService) UpdateAssignment(updateAssignmentRequest *model.UpdateAssignmentRequest, id int) (*model.Assignment, error) {
	assignment, err := a.assignmentRepository.GetByID(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	updateAssignmentRequest.UpdateAssignment(assignment)
	err = a.assignmentRepository.Update(assignment)
	return assignment, err
}

func (a *assignmentService) GetAssignmentById(assignmentId int) (model.AssignmentPublic, error) {
	assignment, err := a.assignmentRepository.GetByID(assignmentId)
	if err != nil {
		logrus.Error(err)
		return model.AssignmentPublic{}, err
	}
	return assignment.ToAssignmentPublic(), nil
}

func (a *assignmentService) DeleteAssignment(assignmentId int) error {
	return a.assignmentRepository.DeleteAssignment(assignmentId)
}
