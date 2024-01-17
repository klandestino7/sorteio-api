package user

import (
	"errors"
	"net/http"
	"sorteio-api/source/dto"
	"sorteio-api/source/models"
	"sorteio-api/source/modules/session"
	"sorteio-api/source/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CONTROLLER
type IUserController interface {
	RequestUserRegistration(c *gin.Context)
	RequestUserAppLogin(c *gin.Context)
	RequestUserPanelLogin(c *gin.Context)
	RequestForgotPassword(c *gin.Context)
	RequestUserFromPhone(c *gin.Context)
}

type UserController struct {
	UserService    IUserService
	SessionService session.ISessionService
}

func InitUserController(UserService IUserService, SessionService session.ISessionService) IUserController {
	return &UserController{
		UserService:    UserService,
		SessionService: SessionService,
	}
}

func (ct *UserController) RequestUserRegistration(c *gin.Context) {
	var json dto.UserRegisterRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)

	userRoles := []models.Role{models.Role{
		Name: "user",
		Role: 0,
	}}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Password: string(password),
		Name:     json.Name,
		Email:    json.Email,
		CPF:      json.Cpf,
		Phone:    json.Phone,
		Blocked:  false,
		Roles:    userRoles,
	}

	_, err := ct.UserService.CreateUser(newUser)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.CreateDetailedErrorDto("database", err))
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"full_messages": []string{"User created successfully"},
	})
}

func (ct *UserController) RequestUserAppLogin(c *gin.Context) {
	var json dto.LoginRequestDto

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	user, err := ct.UserService.GetUserFromEmail(json.Email)

	if err != nil {
		c.JSON(http.StatusForbidden, dto.CreateDetailedErrorDto("login_error", err))
		return
	}

	if user.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, dto.CreateDetailedErrorDto("login", errors.New("invalid credentials")))
		return
	}

	loginData := dto.CreateLoginSuccessful(&user)

	c.JSON(http.StatusOK, loginData)
}

func (ct *UserController) RequestUserPanelLogin(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "content-type, Accept, authorization")

	var json dto.LoginRequestDto

	//validate the request body
	if err := c.BindJSON(&json); err != nil {
		panic(err)
	}

	user, err := ct.UserService.GetUserFromEmail(json.Email)

	if err != nil {
		c.JSON(http.StatusForbidden, dto.CreateDetailedErrorDto("login_error", err))
		return
	}

	if user.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, dto.CreateDetailedErrorDto("login", errors.New("invalid credentials")))
		return
	}

	if !user.IsReseller() {
		if !user.IsAdmin() {
			c.JSON(http.StatusUnauthorized, dto.CreateDetailedErrorDto("login", errors.New("don't have permission")))
			return
		}
	}

	loginData := dto.CreateLoginSuccessful(&user)

	c.JSON(http.StatusOK, loginData)
}

func (ct *UserController) RequestForgotPassword(c *gin.Context) {

	c.JSON(http.StatusOK, dto.CreateSuccessWithMessageDto("GotMessage"))
}

func (ct *UserController) RequestUserFromPhone(c *gin.Context) {
	phone := c.Param("phone")

	if phone == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user Invalid"})
		return
	}

	utils.DebugPrint(" phone :: ", phone)

	response, _ := ct.UserService.GetUserFromCredential("phone", phone)

	user := dto.CreateUserResponse(&response)

	utils.DebugPrint(" user :: ", user)

	c.IndentedJSON(http.StatusOK, user)
}

func (ct *UserController) RequestUserIsAdmin(c *gin.Context) {
	
}

// func (ct *UserController) RequestUserFromToken(c *gin.Context) {
// 	token := c.Param("token")

// 	user, err := ct.SessionService.ReturnUserFromSessionToken(token)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, dto.CreateDetailedErrorDto("not-found-user", errors.New("cant find user")))
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"user": user})
// }
