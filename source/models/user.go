package models

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type RolesEnum uint8

const (
	eRoleUser RolesEnum = iota
	eRoleReseller
	eRoleMod
	eRoleAdmin
)

type Role struct {
	//	ID   		primitive.ObjectID 		`bson:"_id,omitempty"`
	Name string    `bson:"name"`
	Role RolesEnum `bson:"role"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	CPF      string             `bson:"cpf"`
	Email    string             `bson:"email"`
	Phone    string             `bson:"phone"`
	Roles    []Role             `bson:"role"`
	Blocked  bool               `bson:"blocked"`
	Password string             `bson:"password,omitempty"`
	Session  primitive.ObjectID `bson:"session,omitempty"`
}

func (user *User) IsReseller() bool {
	for _, role := range user.Roles {
		if role.Role == eRoleReseller {
			return true
		}
	}
	return false
}

func (user *User) IsAdmin() bool {
	for _, role := range user.Roles {
		if role.Role == eRoleAdmin {
			return true
		}
	}
	return false
}

func (user *User) IsNotAdmin() bool {
	return !user.IsAdmin()
}

// What's bcrypt? https://en.wikipedia.org/wiki/Bcrypt
// Golang bcrypt doc: https://godoc.org/golang.org/x/crypto/bcrypt
// You can change the value in bcrypt.DefaultCost to adjust the security index.
//
//	err := userModel.setPassword("password0")
func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)
	return nil
}

// Database will only save the hashed string, you should check it by util function.
//
//	if err := serModel.checkPassword("password0"); err != nil { password error }
func (u *User) IsValidPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

type TokenCustomClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func getSecretKey() string {
	var secret = os.Getenv("JWT_SECRET_WORD")

	if secret == "" {
		secret = "secret"
	}
	return secret
}

// Generate JWT token associated to this user
func (user *User) GenerateJwtToken() string {
	claims := &TokenCustomClaims{
		user.ID.String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
			Issuer:    "iriffa.com",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		panic(err)
	}

	return t
}
