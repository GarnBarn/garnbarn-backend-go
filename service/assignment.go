package service

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
)

type AssignmentService interface {
	CreateAssignment(assignment *model.Assignment) error
	DeleteAssignment(assignmentId int) error
	GetAllAssignment(fromPresent bool) ([]model.Assignment, error)
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

<<<<<<< HEAD
func (a *assignmentService) DeleteAssignment(assignmentId int) error {
	return a.assignmentRepository.DeleteAssignment(assignmentId)
=======
func (a *assignmentService) GetAllAssignment(fromPresent bool) ([]model.Assignment, error) {
	return a.assignmentRepository.GetAllAssignment(fromPresent)
>>>>>>> 9539c11 ([GB-9] Implement Get All Assignment API (#8))
}
