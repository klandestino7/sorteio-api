package session

import (
	"context"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/models"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

// SERVICE
type ISessionService interface {
	CreateSession(sessionData models.Session) (interface{}, bool)
	GetASession(sessionId string) (models.Session, bool)
	GetSessionFromToken(token string) (models.Session, error)
	// ReturnUserFromSessionToken(token string) (models.User, error)
}

type SessionService struct {
	SessionRepository ISessionRepository
	Validate          *validator.Validate
	// UserService       user.IUserService
}

func InitSessionService(SessionRepository ISessionRepository, validate *validator.Validate) ISessionService {
	return &SessionService{
		SessionRepository: SessionRepository,
		Validate:          validate,
		// UserService:       UserService,
	}
}

func (session *SessionService) CreateSession(sessionHandleData models.Session) (interface{}, bool) {
	if validationErr := session.Validate.Struct(&sessionHandleData); validationErr != nil {
		panic(validationErr)
	}

	sessionResult, _ := session.SessionRepository.Create(sessionHandleData)

	return sessionResult, sessionResult != nil
}

func (session *SessionService) GetASession(sessionId string) (models.Session, bool) {
	sessionSingle := session.SessionRepository.FindOne(sessionId)
	return sessionSingle, !sessionSingle.ID.IsZero()
}

func (session *SessionService) GetSessionFromToken(token string) (models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("session")
	defer cancel()

	var sessionSingle models.Session

	err := coll.FindOne(ctx, bson.M{"token": token}).Decode(&sessionSingle)

	return sessionSingle, err
}

// func (session *SessionService) ReturnUserFromSessionToken(token string) (models.User, error) {
// 	sessionSingle, err := session.GetSessionFromToken(token)

// 	if err != nil {
// 		panic(err)
// 	}

// 	userId := sessionSingle.UserId

// 	user, status := session.UserService.GetUserFromId(userId)

// 	if !status {
// 		panic("err")
// 	}

// 	return user, err
// }
