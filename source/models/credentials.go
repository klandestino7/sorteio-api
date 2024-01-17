package models

type Credentials struct {
	Email string `json:"email"`
	CPF   string `json:"cpf"`
	Phone string `json:"phone"`
}
