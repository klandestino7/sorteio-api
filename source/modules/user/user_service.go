package user

import (
	"errors"
	"fmt"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/models"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SERVICE
type IUserService interface {
	CreateUser(userHandleData models.User) (interface{}, error)
	GetUserFromCredential(credential string, value string) (models.User, error)
	GetUserFromEmail(email string) (models.User, error)
	GetUserFromId(userId primitive.ObjectID) (models.User, bool)
	EmailHasTaken(email string) bool
	CPFHasTaken(cpf string) bool
	PhoneNumberHasTaken(phone string) bool
	GetUserFromCredentials(credentials models.Credentials) models.User
}

type UserService struct {
	UserRepository IUserRepository
	Validate       *validator.Validate
}

func InitUserService(UserRepository IUserRepository, validate *validator.Validate) IUserService {
	return &UserService{
		UserRepository: UserRepository,
		Validate:       validate,
	}
}

func (s *UserService) CreateUser(userHandleData models.User) (interface{}, error) {
	if validationErr := s.Validate.Struct(&userHandleData); validationErr != nil {
		panic(validationErr)
	}

	var credentialsUsed = []bool{
		s.EmailHasTaken(userHandleData.Email),
		s.CPFHasTaken(userHandleData.CPF),
		s.PhoneNumberHasTaken(userHandleData.Phone),
	}

	var credentialsFromNumber = []string{
		"email",
		"cpf",
		"phone",
	}

	var errs error = nil

	for i, value := range credentialsUsed {
		if value {
			errs = errors.New(fmt.Sprintf("Credential %s has taken", credentialsFromNumber[i]))
		}
	}

	userResult, _ := s.UserRepository.Create(userHandleData)

	return userResult, errs
}

func (s *UserService) GetUserFromCredential(credential string, value string) (models.User, error) {
	filter := bson.M{fmt.Sprintf("%s", credential): value}

	result, _ := DBConnection.FindADocument("user", filter)

	var user models.User
	err := result.Decode(&user)

	return user, err
}

func (s *UserService) GetUserFromEmail(email string) (models.User, error) {
	filter := bson.M{"email": email}
	result, _ := DBConnection.FindADocument("user", filter)

	var user models.User
	err := result.Decode(&user)

	return user, err
}

func (user *UserService) GetUserFromId(userId primitive.ObjectID) (models.User, bool) {
	userSingle := user.UserRepository.FindOne(userId.String())
	return userSingle, !userSingle.ID.IsZero()
}

func (user *UserService) EmailHasTaken(email string) bool {
	userSingle := user.UserRepository.FindOneWithFilter(bson.M{"email": email})

	if userSingle.Email == email {
		return true
	}

	return false
}

func (user *UserService) CPFHasTaken(cpf string) bool {
	userSingle := user.UserRepository.FindOneWithFilter(bson.M{"cpf": cpf})

	if userSingle.CPF == cpf {
		return true
	}

	return false
}

func (user *UserService) PhoneNumberHasTaken(phone string) bool {
	userSingle := user.UserRepository.FindOneWithFilter(bson.M{"phone": phone})

	if userSingle.Phone == phone {
		return true
	}

	return false
}

func (user *UserService) GetUserFromCredentials(credentials models.Credentials) models.User {

	var credentialsUsed = []bool{
		user.EmailHasTaken(credentials.Email),
		user.CPFHasTaken(credentials.CPF),
		user.PhoneNumberHasTaken(credentials.Phone),
	}

	var credentialsFromNumber = []string{
		"email",
		"cpf",
		"phone",
	}

	valueFromCredential := map[string]string{
		"email": credentials.Email,
		"cpf":   credentials.CPF,
		"phone": credentials.Phone,
	}

	for i, value := range credentialsUsed {
		if value {
			var credential string = credentialsFromNumber[i]
			var value string = valueFromCredential[credential]

			user, err := user.GetUserFromCredential(credential, value)

			if err != nil {
				panic(err)
			}

			return user
		}
	}

	return models.User{}
}
