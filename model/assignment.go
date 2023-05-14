package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/GarnBarn/garnbarn-backend-go/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Assignment struct {
	gorm.Model
	Name         string
	Author       string
	Description  string
	ReminderTime string
	DueDate      int
	TagID        int
	Tag          *Tag
}

func (a *Assignment) BeforeSave(tx *gorm.DB) (err error) {
	// Encrypt the data before saving into the database
	key := tx.Statement.Context.Value(config.AssignmentEncryptionContextKey).(string)
	a.Name, err = RemainOrEncrypt(a.Name, key)
	if err != nil {
		logrus.Error("Encrypt Data Error: ", err)
		return err
	}

	a.Description, err = RemainOrEncrypt(a.Description, key)
	if err != nil {
		logrus.Error("Encrypt Data Error: ", err)
		return err
	}

	a.ReminderTime, err = RemainOrEncrypt(a.ReminderTime, key)
	if err != nil {
		logrus.Error("Encrypt Data Error: ", err)
		return err
	}

	return nil
}

func (a *Assignment) AfterFind(tx *gorm.DB) (err error) {
	// Decrypt the data.

	key := tx.Statement.Context.Value(config.AssignmentEncryptionContextKey).(string)
	a.Name, err = DecryptAES(key, a.Name)
	if err != nil {
		logrus.Error("Decrypt Data Error: ", err)
		return err
	}

	a.Description, err = RemainOrDecrypt(a.Description, key)
	if err != nil {
		logrus.Error("Decrypt Data Error: ", err)
		return err
	}

	a.ReminderTime, err = RemainOrDecrypt(a.ReminderTime, key)
	if err != nil {
		logrus.Error("Decrypt Data Error: ", err)
		return err
	}

	return nil
}

func (a *Assignment) ToAssignmentPublic() AssignmentPublic {
	reminderTime := strings.Split(a.ReminderTime, ",")
	reminterTimeInt := []int{}

	for _, item := range reminderTime {
		result, err := strconv.Atoi(item)
		if err != nil {
			logrus.Warn("Can't convert the result to int: ", item, " for ", a.ID)
			continue
		}
		reminterTimeInt = append(reminterTimeInt, result)
	}

	assignmentResult := AssignmentPublic{
		ID:           fmt.Sprint(a.ID),
		Name:         a.Name,
		Author:       a.Author,
		Description:  a.Description,
		DueDate:      a.DueDate,
		Tag:          nil,
		ReminderTime: reminterTimeInt,
	}

	if a.Tag != nil {
		tagPublicResult := a.Tag.ToTagPublic(true)
		assignmentResult.Tag = &tagPublicResult
	}
	return assignmentResult
}

type AssignmentPublic struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Author       string     `json:"author"`
	Description  string     `json:"description,omitempty"`
	DueDate      int        `json:"dueDate"`
	Tag          *TagPublic `json:"tag"`
	ReminderTime []int      `json:"reminderTime"`
}

type AssignmentRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	DueDate      int    `json:"dueDate"`
	TagId        string `json:"tagId"`
	ReminderTime []int  `json:"reminderTime,omitempty" validate:"max=3,omitempty"`
}

func (ar *AssignmentRequest) ToAssignment(author string) Assignment {
	reminderTimeByte, _ := json.Marshal(ar.ReminderTime)
	reminderTimeString := strings.Trim(string(reminderTimeByte), "[]")

	tagIdInt, _ := strconv.Atoi(ar.TagId)

	tag := Tag{}
	tag.ID = uint(tagIdInt)

	return Assignment{
		Name:         ar.Name,
		Author:       author,
		Description:  ar.Description,
		ReminderTime: reminderTimeString,
		DueDate:      ar.DueDate,
		TagID:        tagIdInt,
		Tag:          &tag,
	}
}

type UpdateAssignmentRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	DueDate      *int    `json:"dueDate,omitempty"`
	TagId        *int    `json:"tagId,omitempty"`
	ReminderTime *[]int  `json:"reminderTime,omitempty" validate:"max=3,omitempty"`
}

func (ur *UpdateAssignmentRequest) UpdateAssignment(assignment *Assignment) {
	if ur.Name != nil {
		assignment.Name = *ur.Name
	}
	if ur.Description != nil {
		assignment.Description = *ur.Description
	}
	if ur.DueDate != nil {
		assignment.DueDate = *ur.DueDate
	}
	if ur.TagId != nil {
		assignment.TagID = *ur.TagId
	}
	if ur.ReminderTime != nil {
		assignment.ReminderTime = convertReminterTimeToString(*ur.ReminderTime)
	}
}
