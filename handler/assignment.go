package handler

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AssignmentHandler struct {
	validate          validator.Validate
	assignmentService service.AssignmentService
	tagService        service.Tag
}

func NewAssignmentHandler(validate validator.Validate, assignmentService service.AssignmentService, tagService service.Tag) AssignmentHandler {
	return AssignmentHandler{
		validate:          validate,
		assignmentService: assignmentService,
		tagService:        tagService,
	}
}

func (a *AssignmentHandler) GetAllAssignment(c *gin.Context) {
	fromPresentString := c.Query("fromPresent")

	logrus.Debug("From Present string: ", fromPresentString)

	fromPresent := true
	if fromPresentString == "" || fromPresentString == "false" {
		fromPresent = false
	}

	logrus.Debug("From Present: ", fromPresent)

	uid := c.Param(UserUidKey)

	assignments, err := a.assignmentService.GetAllAssignment(uid, fromPresent)

	assignmentPublic := []model.AssignmentPublic{}

	for _, item := range assignments {
		assignmentPublic = append(assignmentPublic, item.ToAssignmentPublic())
	}

	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in server."})
		return
	}

	c.JSON(http.StatusOK, model.BulkResponse[model.AssignmentPublic]{
		Count:    len(assignmentPublic),
		Previous: nil,
		Next:     nil,
		Results:  assignmentPublic,
	})
}

func (a *AssignmentHandler) CreateAssignment(c *gin.Context) {
	var assignmentRequest model.AssignmentRequest

	err := c.ShouldBindJSON(&assignmentRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = a.validate.Struct(assignmentRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	assignment := assignmentRequest.ToAssignment(c.Param(UserUidKey))

	if err := a.assignmentService.CreateAssignment(&assignment); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something happen in server.",
		})
		return
	}

	assignmentPublic := assignment.ToAssignmentPublic()

	c.JSON(http.StatusCreated, assignmentPublic)

}

func (a *AssignmentHandler) UpdateAssignment(c *gin.Context) {
	assignmentIdString, ok := c.Params.Get("assignmentId")
	if !ok {
		logrus.Warn("Can't get assignmentId from parameters")
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	// Check if the tagId is int parsable
	assignmentId, err := strconv.Atoi(assignmentIdString)
	if err != nil {
		logrus.Warn("Can't convert assignmentId to int: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	// Bind the request body.
	var updateAssignmentRequest model.UpdateAssignmentRequest
	err = c.ShouldBindJSON(&updateAssignmentRequest)
	if err != nil {
		logrus.Warn("Can't bind request body to model: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	err = a.validate.Struct(updateAssignmentRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//Check if tagId is existed
	updateTagIdRequest := updateAssignmentRequest.TagId
	updateTagIdRequestInt, err := strconv.Atoi(*updateTagIdRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if updateTagIdRequest != nil && !a.tagService.IsTagExist(updateTagIdRequestInt) {
		logrus.Warn("Tag id is not exist")
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	updatedAssignment, err := a.assignmentService.UpdateAssignment(&updateAssignmentRequest, assignmentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in the server"})
		return
	}

	publicAssignment := updatedAssignment.ToAssignmentPublic()
	c.JSON(http.StatusOK, publicAssignment)
}

func (a *AssignmentHandler) DeleteAssignment(c *gin.Context) {
	assignmentIdString, ok := c.Params.Get("Id")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	assignmentId, err := strconv.Atoi(assignmentIdString)
	err = a.assignmentService.DeleteAssignment(assignmentId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func (a *AssignmentHandler) GetAssignmentById(c *gin.Context) {
	assignmentIdStr, ok := c.Params.Get("assignmentId")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	assignmentId, err := strconv.Atoi(assignmentIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	publicAssignment, err := a.assignmentService.GetAssignmentById(assignmentId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "assignment id not found"})
		return
	}
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in server."})
		return
	}
	c.JSON(http.StatusOK, publicAssignment)
}
