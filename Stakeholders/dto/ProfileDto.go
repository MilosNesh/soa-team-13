package dto

type ProfileDto struct {
	Name           *string `json:"name,omitempty"`
	Surname        *string `json:"surname,omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	Bio            *string `json:"bio,omitempty"`
	Motto          *string `json:"motto,omitempty"`
}
