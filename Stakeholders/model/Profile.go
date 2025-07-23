package model

type Profile struct {
	AccountId      string `gorm:"type:uuid;primaryKey"` //primary i foreign key
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
	Motto          string `json:"motto"`
}
