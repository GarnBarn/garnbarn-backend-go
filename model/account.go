package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Uid     string
	Line    string
	Consent bool
}

func (a *Account) ToAccountPublic(displayName string, profileImage string) AccountPublic {
	return AccountPublic{
		Uid:          a.Uid,
		DisplayName:  displayName,
		ProfileImage: profileImage,
		Platform: &AccountPlatform{
			Line: a.Line,
		},
		Consent: a.Consent,
	}
}

type AccountPublic struct {
	Uid          string           `json:"uid"`
	DisplayName  string           `json:"displayName"`
	ProfileImage string           `json:"profileImage"`
	Platform     *AccountPlatform `json:"platform"`
	Consent      bool             `json:"consent"`
}

type AccountPlatform struct {
	Line string `json:"line"`
}

type UpdateAccountRequest struct {
	Uid     string `json:"uid"`
	Consent bool   `json:"consent"`
}

type CheckCompromisedPasswordRequest struct {
	HashedPassword string `json:"hashedPassword"`
}
