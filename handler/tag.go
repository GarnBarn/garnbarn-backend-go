package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Tag struct {
	validate   validator.Validate
	tagService service.Tag
}

var (
	ErrGinBadRequestBody = gin.H{"message": "bad request body."}
)

func NewTagHandler(validate validator.Validate, tagService service.Tag) Tag {
	return Tag{
		validate:   validate,
		tagService: tagService,
	}
}

func (t *Tag) CreateTag(c *gin.Context) {

	var tagRequest model.CreateTagRequest

	err := c.ShouldBind(&tagRequest)
	if err != nil {
		logrus.Warn("Requets Body binding error: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	err = t.validate.Struct(tagRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// TODO: Change this to the actual user from the authentication header
	tag := tagRequest.ToTag("test")

	err = t.tagService.CreateTag(&tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in the server"})
		return
	}

	tagPublic := tag.ToTagPublic(false)
	c.JSON(http.StatusOK, tagPublic)
}

func (t *Tag) UpdateTag(c *gin.Context) {
	tagIdString, ok := c.Params.Get("tagId")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	// Check if the tagId is int parsable
	tagId, err := strconv.Atoi(tagIdString)
	if err != nil {
		logrus.Warn("Can't convert tagId to int: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	// Bind the request body.
	var updateTagRequest model.UpdateTagRequest
	err = c.ShouldBind(&updateTagRequest)
	if err != nil {
		logrus.Warn("Can't bind request body to model: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	err = t.validate.Struct(updateTagRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tag, err := t.tagService.UpdateTag(tagId, &updateTagRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in the server"})
		return
	}

	tagPublic := tag.ToTagPublic(true)
	c.JSON(http.StatusOK, tagPublic)
}

func (t *Tag) GetTagById(c *gin.Context) {
	tagIdStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	publicTag, err := t.tagService.GetTagById(tagId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrGinBadRequestBody)
		return
	}

	c.JSON(http.StatusOK, publicTag)
}
