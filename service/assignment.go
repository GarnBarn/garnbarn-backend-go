package service

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/sirupsen/logrus"
)

type AssignmentService interface {
	CreateAssignment(assignment *model.Assignment) error
	GetAllAssignment(fromPresent bool) ([]model.Assignment, error)
	UpdateAssignment(updateAssignmentRequest *model.UpdateAssignmentRequest, id int) (*model.Assignment, error)
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

func (a *assignmentService) GetAllAssignment(fromPresent bool) ([]model.Assignment, error) {
	return a.assignmentRepository.GetAllAssignment(fromPresent)
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
