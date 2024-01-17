package routes

import (
	"sorteio-api/source/modules/winner"

	"github.com/gin-gonic/gin"
)

func WinnerRoute(router *gin.Engine, winner winner.IWinnerController) {
	router.GET("/winners", winner.RequestAllWinners)
}
