package models

//Login credential
type LoginCredentials struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}