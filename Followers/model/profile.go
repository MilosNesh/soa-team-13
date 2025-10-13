package model

type Profile struct {
	Id       string `json:"id"`
	Username string `json:"username,omitempty"`
}
