package models

type User struct {
	Id          uint   `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	Is_verified bool   `json:"is_Verified,omitempty"`
}
