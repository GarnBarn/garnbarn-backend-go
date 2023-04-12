package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
	TagId        int
}

func (a *Assignment) ToAssignmentResponse() AssignmentResponse {
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

	return AssignmentResponse{
		ID:           fmt.Sprint(a.ID),
		Name:         a.Name,
		Author:       a.Author,
		Description:  a.Description,
		DueDate:      a.DueDate,
		TagId:        a.TagId,
		ReminderTime: reminterTimeInt,
	}
}

type AssignmentResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Author       string `json:"author"`
	Description  string `json:"description,omitempty"`
	DueDate      int    `json:"dueDate"`
	TagId        int    `json:"tagId"`
	ReminderTime []int  `json:"reminderTime"`
}

type AssignmentRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	DueDate      int    `json:"dueDate"`
	TagId        int    `json:"tagId" validate:"min=1"`
	ReminderTime []int  `json:"reminderTime,omitempty" validate:"len=3,omitempty"`
}

func (ar *AssignmentRequest) ToAssignment(author string) Assignment {
	reminderTimeByte, _ := json.Marshal(ar.ReminderTime)
	reminderTimeString := strings.Trim(string(reminderTimeByte), "[]")

	return Assignment{
		Name:         ar.Name,
		Author:       author,
		Description:  ar.Description,
		ReminderTime: reminderTimeString,
		DueDate:      ar.DueDate,
		TagId:        ar.TagId,
	}
}
