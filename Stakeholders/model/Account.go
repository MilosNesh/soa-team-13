package model

type Account struct {
	Id       string `gorm:"type:uuid;primaryKey" json:"id"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Role     string `json:"role"`
	Blocked  bool   `json:"blocked" gorm:"not null;default:false;index"`

	Profile Profile `gorm:"constraint:OnDelete:CASCADE"`
}

func (a *Account) IsValid() bool {
	if a.Username == "" || a.Password == "" || a.Email == "" || a.Role == "" {
		return false
	}
	return true
}
