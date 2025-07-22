package model

type Account struct {
	Username string `json:"username" gorm:"primaryKey;not null"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Role     string `json:"role"`

	Profile Profile `gorm:"foreignKey:Username;references:Username;constraint:OnDelete:CASCADE"`
}

func (a *Account) IsValid() bool {
	if a.Username == "" || a.Password == "" || a.Email == "" || a.Role == "" {
		return false
	}
	return true
}
