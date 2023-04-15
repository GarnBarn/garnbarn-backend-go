package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Uid  string
	Line string
}

func (a *Account) ToAccountPublic(displayName string, profileImage string) AccountPublic {
	return AccountPublic{
		Uid:          a.Uid,
		DisplayName:  displayName,
		ProfileImage: profileImage,
		Platform: &AccountPlatform{
			Line: a.Line,
		},
	}
}

type AccountPublic struct {
	Uid          string           `json:"uid"`
	DisplayName  string           `json:"displayName"`
	ProfileImage string           `json:"profileImage"`
	Platform     *AccountPlatform `json:"platform"`
}

type AccountPlatform struct {
	Line string `json:"line"`
}
