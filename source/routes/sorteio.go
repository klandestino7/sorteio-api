package routes

import (
	"sorteio-api/source/modules/sorteio"

	"github.com/gin-gonic/gin"
)

func SorteioRoute(router *gin.Engine, sorteio sorteio.ISorteioController) {
	router.GET("/sorteios/", sorteio.RequestAllSorteios)
	router.GET("/sorteios/:sorteioId", sorteio.RequestSorteioFromId)
	router.POST("/sorteios/winner/generate/:sorteioId", sorteio.RequestGenerateWinner)
}

func SorteioProtectedRoute(router *gin.RouterGroup, sorteio sorteio.ISorteioController) {
	router.GET("/sorteios/full", sorteio.RequestFullSorteios)
	router.POST("/sorteios/create", sorteio.RequestCreateNewSorteio)
	router.POST("/sorteios/update", sorteio.RequestUpdateSorteio)
}
