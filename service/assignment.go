package service

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
)

type AssignmentService interface {
	CreateAssignment(assignment *model.Assignment) error
	GetAllAssignment() ([]model.Assignment, error)
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

func (a *assignmentService) GetAllAssignment() ([]model.Assignment, error) {
	return a.assignmentRepository.GetAllAssignment()
}
