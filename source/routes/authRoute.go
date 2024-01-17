package routes

import (
	auth "sorteio-api/source/modules/auth_token"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine, auth auth.IAuthTokenController) {
	router.POST("/auth/user/create", auth.CreateSessionToken)
}
