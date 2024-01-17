package authToken

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

// SERVICE
type IAuthTokenService interface {
	GenerateToken(userId string) string
	ValidateToken(token string) (*jwt.Token, error)
	GetTokenUserId(claims jwt.MapClaims) string
	CheckIsWebhookAutenticationToken(token string) bool
}

type AuthTokenService struct {
	AuthTokenRepository IAuthTokenRepository
	Validate            *validator.Validate
	secretKey           string
	issure              string
}

type TokenCustomClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func InitAuthTokenService(authTokenRepository IAuthTokenRepository, validate *validator.Validate) IAuthTokenService {
	return &AuthTokenService{
		AuthTokenRepository: authTokenRepository,
		Validate:            validate,
		secretKey:           getSecretKey(),
		issure:              "iriffa.com",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET_WORD")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (s *AuthTokenService) GenerateToken(userId string) string {
	claims := &TokenCustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
			Issuer:    s.issure,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		panic(err)
	}

	return t
}

func (s *AuthTokenService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if res, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			fmt.Println("RES :: ", res)
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(getSecretKey()), nil
	})
}

func (s *AuthTokenService) GetTokenUserId(claims jwt.MapClaims) string {
	fmt.Println(claims["user_id"])
	return fmt.Sprint(claims["user_id"])
}

func (s *AuthTokenService) CheckIsWebhookAutenticationToken(token string) bool {
	authKey := os.Getenv("AUTH_WEBHOOK_TOKEN")
	return authKey == token
}
