package model

type Profile struct {
	Username       string `json:"username" gorm:"primaryKey;not null"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
	Motto          string `json:"motto"`
}
