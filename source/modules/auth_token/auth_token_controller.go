package authToken

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CONTROLLER
type IAuthTokenController interface {
	RequestGetWebhookAuthentication(c *gin.Context)
	CreateSessionToken(c *gin.Context)
}

type AuthTokenController struct {
	AuthTokenService IAuthTokenService
}

func InitAuthTokenController(authTokenService IAuthTokenService) IAuthTokenController {
	return &AuthTokenController{
		AuthTokenService: authTokenService,
	}
}

func (ct *AuthTokenController) RequestGetWebhookAuthentication(c *gin.Context) {
	tokenId := c.Param("token")

	if tokenId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "transactionId Invalid"})
		return
	}

	response := ct.AuthTokenService.CheckIsWebhookAutenticationToken(tokenId)

	if !response {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": response})
}

func (ct *AuthTokenController) CreateSessionToken(c *gin.Context) {
	sessionToken := ct.AuthTokenService.GenerateToken(c.ClientIP())

	if sessionToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": sessionToken})
}
