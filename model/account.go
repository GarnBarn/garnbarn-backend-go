package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Uid  string
	Line string
}

type AccountPublic struct {
	Uid         string          `json:"uid"`
	DisplayName string          `json:"displayName"`
	Platform    AccountPlatform `json:"platform"`
}

type AccountPlatform struct {
	Line string `json:"line"`
}
