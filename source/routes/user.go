package routes

import (
	"sorteio-api/source/modules/user"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine, user user.IUserController) {
	router.POST("/users/login-panel", user.RequestUserPanelLogin)
	router.POST("/users/forgot-password", user.RequestForgotPassword)

	// router.GET("/users/get-from-token/:token", user.RequestUserFromToken)
	router.GET("/users/get-from-phone/:phone", user.RequestUserFromPhone)
}

func UserProtectedRoute(router *gin.RouterGroup, user user.IUserController) {
	router.POST("/users/create", user.RequestUserRegistration)
}
