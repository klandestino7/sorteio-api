package routes

import (
	"sorteio-api/source/modules/order"

	"github.com/gin-gonic/gin"
)

func OrderRoute(router *gin.Engine, order order.IOrderController) {
	router.GET("/orders", order.RequestAllOrders)
	router.GET("/orders/get-id/:orderId", order.RequestOrderFromId)
	router.GET("/orders/get-objId/:orderId", order.RequestOrderFromObjId)
	router.GET("/orders/get-transaction/:transactionId", order.RequestOrderStatusFromTransactionId)
	router.GET("/orders/from-sorteio", order.RequestOrdersFromSorteioId)

	router.POST("/orders/create", order.RequestCreateOrder)
	router.POST("/orders/from-credential", order.RequestOrdersFromCredential)
}

func OrderProtectedRoute(router *gin.RouterGroup, order order.IOrderController) {
	router.GET("/orders/paids", order.RequestOrdersPaid)
}
