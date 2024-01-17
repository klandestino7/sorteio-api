package session

import "github.com/gin-gonic/gin"

// CONTROLLER
type ISessionController interface {
	// RequestSessionFromIndex()
}

type SessionController struct {
	SessionService ISessionService
}

func InitSessionController(SessionService ISessionService) ISessionController {
	return &SessionController{
		SessionService: SessionService,
	}
}

func RequestStartNewSession(c *gin.Context) {

}
