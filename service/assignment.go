package service

import (
	"net/http"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
	"github.com/gin-gonic/gin"
)

type AssignmentService interface {
	CreateAssignment(c gin.Context)
}

type assignmentService struct {
	assignmentRepository repository.AssignmentRepository
}

func NewAssignmentService(assignmentRepository repository.AssignmentRepository) AssignmentService {
	return &assignmentService{
		assignmentRepository: assignmentRepository,
	}
}

func (a *assignmentService) CreateAssignment(c gin.Context) {
	var assignment model.Assignment

	// Check the conditional operator (?:) later.
	// add validate
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result := a.assignmentRepository.CreateAssignment(assignment); result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result,
		})
		return
	}

	c.JSON(http.StatusCreated, &assignment)
}
