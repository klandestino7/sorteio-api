package dto

import (
	"sorteio-api/source/models"
)

type UserRegisterRequestDto struct {
	Name                 string `form:"name" json:"name" xml:"name" binding:"required"`
	Cpf                  string `form:"cpf" json:"cpf" xml:"cpf" binding:"required"`
	Email                string `form:"email" json:"email" xml:"email" binding:"required"`
	Phone                string `form:"phone" json:"phone" xml:"phone" binding:"required"`
	Password             string `form:"password" json:"password" xml:"password" binding:"required"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" xml:"password-confirmation" binding:"required"`
}

type LoginRequestDto struct {
	//	 Username string `form:"username" json:"username" xml:"username" binding:"exists,username"`
	//	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Email    string `form:"email" json:"email" xml:"email" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	//	Password string `form:"password" json:"password" binding:"exists,min=8,max=255"`

	//	userModel models.User `json:"-"`
}

func CreateLoginSuccessful(user *models.User) map[string]interface{} {
	var roles = make([]string, len(user.Roles))

	for i := 0; i < len(user.Roles); i++ {
		roles[i] = user.Roles[i].Name
	}

	jwtToken := user.GenerateJwtToken()

	handleData := map[string]interface{}{
		// "success": true,
		"token": jwtToken,
		"user": map[string]interface{}{
			// "email": user.Email,
			"id":    user.ID,
			"roles": roles,
		},
	}

	return handleData
}

func CreateUserResponse(order *models.User) map[string]interface{} {
	handleData := map[string]interface{}{
		"id":      order.ID,
		"name":    order.Name,
		"cpf":     order.CPF,
		"email":   order.Email,
		"phone":   order.Phone,
		"roles":   order.Roles,
		"blocked": order.Blocked,
	}

	return handleData
}
