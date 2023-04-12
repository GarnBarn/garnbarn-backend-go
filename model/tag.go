package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name          string
	Author        string
	Color         string
	ReminderTime  string
	Subscriber    string
	SecretKeyTotp string
}

func (t *Tag) ToTagPublic() TagPublic {
	reminderTime := strings.Split(t.ReminderTime, ",")
	reminterTimeInt := []int{}

	for _, item := range reminderTime {
		result, err := strconv.Atoi(item)
		if err != nil {
			logrus.Warn("Can't convert the result to int: ", item, " for ", t.ID)
			continue
		}
		reminterTimeInt = append(reminterTimeInt, result)
	}

	return TagPublic{
		ID:            fmt.Sprint(t.ID),
		Name:          t.Name,
		Author:        t.Author,
		Color:         t.Color,
		ReminderTime:  reminterTimeInt,
		Subscriber:    strings.Split(t.Subscriber, ","),
		SecretKeyTotp: t.SecretKeyTotp,
	}
}

type TagPublic struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Author        string   `json:"author"`
	Color         string   `json:"color"`
	ReminderTime  []int    `json:"reminderTime"`
	Subscriber    []string `json:"subscribe"`
	SecretKeyTotp string   `json:"secretKeyTotp,omitempty"`
}

type CreateTagRequest struct {
	Name         string   `json:"name" validate:"required"`
	Color        string   `json:"color"`
	ReminderTime []int    `json:"reminderTime,omitempty" validate:"len=3,omitempty"`
	Subscriber   []string `json:"subscribe"`
}

func (ct *CreateTagRequest) ToTag(author string) Tag {
	reminderTimeByte, _ := json.Marshal(ct.ReminderTime)
	reminderTimeString := strings.Trim(string(reminderTimeByte), "[]")

	return Tag{
		Name:         ct.Name,
		Author:       author,
		Color:        ct.Color,
		ReminderTime: reminderTimeString,
		Subscriber:   strings.Join(ct.Subscriber, ","),
	}
}

type UpdateTagRequest struct {
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	ReminderTime []int  `json:"reminderTime,omitempty" validate:"len=3,omitempty"`
}
