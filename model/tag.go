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

func convertReminterTimeToString(reminterTime []int) string {
	reminderTimeByte, _ := json.Marshal(reminterTime)
	return strings.Trim(string(reminderTimeByte), "[]")
}

func (t *Tag) ToTagPublic(maskSecretKey bool) TagPublic {
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

	secretKey := ""
	if !maskSecretKey {
		secretKey = t.SecretKeyTotp
	}

	return TagPublic{
		ID:            fmt.Sprint(t.ID),
		Name:          t.Name,
		Author:        t.Author,
		Color:         t.Color,
		ReminderTime:  reminterTimeInt,
		Subscriber:    strings.Split(t.Subscriber, ","),
		SecretKeyTotp: secretKey,
	}
}

type TagPublic struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Author        string   `json:"author"`
	Color         string   `json:"color"`
	ReminderTime  []int    `json:"reminderTime"`
	Subscriber    []string `json:"subscriber"`
	SecretKeyTotp string   `json:"secretKeyTotp,omitempty"`
}

type CreateTagRequest struct {
	Name         string   `json:"name" validate:"required"`
	Color        string   `json:"color"`
	ReminderTime []int    `json:"reminderTime,omitempty" validate:"omitempty,max=3"`
	Subscriber   []string `json:"subscriber"`
}

func (ct *CreateTagRequest) ToTag(author string) Tag {
	return Tag{
		Name:         ct.Name,
		Author:       author,
		Color:        ct.Color,
		ReminderTime: convertReminterTimeToString(ct.ReminderTime),
		Subscriber:   strings.Join(ct.Subscriber, ","),
	}
}

type UpdateTagRequest struct {
	Name         *string   `json:"name,omitempty"`
	Color        *string   `json:"color,omitempty"`
	ReminderTime *[]int    `json:"reminderTime,omitempty" validate:"omitempty,max=3"`
	Subscriber   *[]string `json:"subscribe"`
}

func (utr *UpdateTagRequest) UpdateTag(tag *Tag) {
	if utr.Name != nil {
		tag.Name = *utr.Name
	}

	if utr.Color != nil {
		tag.Color = *utr.Color
	}

	if utr.ReminderTime != nil {
		tag.ReminderTime = convertReminterTimeToString(*utr.ReminderTime)
	}

	if utr.Subscriber != nil {
		tag.Subscriber = strings.Join(*utr.Subscriber, ",")
	}
}
