package dto

type ProfileDto struct {
	AccountId      string  `json:"account_id"`
	Name           *string `json:"name,omitempty"`
	Surname        *string `json:"surname,omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	Bio            *string `json:"bio,omitempty"`
	Motto          *string `json:"motto,omitempty"`
}
