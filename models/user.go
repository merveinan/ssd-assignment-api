package models

// User represents the structure of a user in the system
type User struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"password123"`
}
