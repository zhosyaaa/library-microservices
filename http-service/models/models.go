package models

type User struct {
	Id               uint   `json:"id,omitempty"`
	Email            string `json:"email,omitempty"`
	Password         string `json:"password,omitempty"`
	Is_verified      bool   `json:"is_Verified,omitempty"`
	VerificationCode string `json:"verificationCode"`
}
type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	PublishedAt string `json:"published_at"`
}
