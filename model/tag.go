package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/GarnBarn/garnbarn-backend-go/config"
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

func RemainOrEncrypt(data string, key string) (string, error) {
	if data == "" {
		return data, nil
	}

	return Encrypt(data, key)
}

func RemainOrDecrypt(data string, key string) (string, error) {
	if data == "" {
		return data, nil
	}

	return DecryptAES(key, data)
}

func (t *Tag) BeforeSave(tx *gorm.DB) (err error) {
	// Encrypt the data before saving into the database
	key := tx.Statement.Context.Value(config.TagEncryptionContextKey).(string)
	t.Name, err = Encrypt(t.Name, key)
	if err != nil {
		logrus.Error("Encrypt Data Error: ", err)
		return err
	}

	t.Color, err = RemainOrEncrypt(t.Color, key)
	if err != nil {
		logrus.Error("Encrypt Data Error: ", err)
		return err
	}

	t.ReminderTime, err = RemainOrEncrypt(t.ReminderTime, key)
	if err != nil {
		logrus.Error("Encrypt Data Error: ", err)
		return err
	}

	t.SecretKeyTotp, err = RemainOrEncrypt(t.SecretKeyTotp, key)
	if err != nil {
		logrus.Error("Encrypt Data Error: ", err)
		return err
	}

	return nil
}

func (t *Tag) AfterFind(tx *gorm.DB) (err error) {
	// Decrypt the data.

	key := tx.Statement.Context.Value(config.TagEncryptionContextKey).(string)
	t.Name, err = RemainOrDecrypt(t.Name, key)
	if err != nil {
		logrus.Error("Tag Decrypt Data Error: ", err)
		return err
	}

	t.Color, err = RemainOrDecrypt(t.Color, key)
	if err != nil {
		logrus.Error("Tag Color Decrypt Data Error: ", err)
		return err
	}

	t.ReminderTime, err = RemainOrDecrypt(t.ReminderTime, key)
	if err != nil {
		logrus.Error("Tag Decrypt Data Error: ", err)
		return err
	}

	t.SecretKeyTotp, err = RemainOrDecrypt(t.SecretKeyTotp, key)
	if err != nil {
		logrus.Error("Tag Decrypt Data Error: ", err)
		return err
	}
	return nil
}

func Encrypt(message string, key string) (encoded string, err error) {
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher([]byte(key))

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return hex.EncodeToString(cipherText), err
}

func DecryptAES(key string, secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := hex.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher([]byte(key))

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
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
