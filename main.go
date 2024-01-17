package main

import (
	"net/http"
	"os"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/middleware"
	"sorteio-api/source/modules"
	efi "sorteio-api/source/resources/efi_sdk"
	"sorteio-api/source/resources/sentry"
	"sorteio-api/source/routes"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Welcome to sorteio Api"})
}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(".env file couldn't be loaded")
	}

	sentry.SentryInitialization()

	DBConnection.StartMongoDBConnection()

	router := gin.Default()
	router.GET("/", home)
	router.Use(Cors())
	router.Use(Options)

	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	authorized := router.Group("/")
	authorized.Use(middleware.Middleware(func(token string, c *gin.Context) bool {

		jwtToken, err := modules.AuthTokenService.ValidateToken(token)

		if err != nil {

			// If token is not valid, set HTTP status code to 401 (Unauthorized)
			// And return false for blocking the request flow
			c.AbortWithStatus(http.StatusForbidden)
			return false
		}

		modules.AuthTokenService.GetTokenUserId(jwtToken.Claims.(jwt.MapClaims))

		// If token is valid, return true to continue the request flow
		return true
	}))
	{
		routes.OrderProtectedRoute(authorized, modules.OrderController)
		routes.SorteioProtectedRoute(authorized, modules.SorteioController)
		routes.UserProtectedRoute(authorized, modules.UserController)

		// routes.CloudRoute(authorized, cloud)
	}

	routes.OrderRoute(router, modules.OrderController)
	routes.UserRoute(router, modules.UserController)

	routes.SorteioRoute(router, modules.SorteioController)
	routes.WinnerRoute(router, modules.WinnerController)
	routes.AuthRoute(router, modules.AuthTokenController)

	router.GET("/ticket", modules.TicketController.TryGetTicket)

	// router.POST("/webhook-config", modules.OrderController.ConfigWebhookRequest)

	authorizedWebhoook := router.Group("/")
	authorizedWebhoook.Use(middleware.Middleware(func(token string, c *gin.Context) bool {
		if !modules.AuthTokenService.CheckIsWebhookAutenticationToken(token) {
			// If token is not valid, set HTTP status code to 401 (Unauthorized)
			// And return false for blocking the request flow
			c.AbortWithStatus(http.StatusForbidden)
			return false
		}
		// If token is valid, return true to continue the request flow
		return true
	}))
	{
		//authorizedWebhoook.POST("/webhook", orderController.WebhookRequest);;
		authorizedWebhoook.POST("/pix", modules.OrderController.PixHandle)
		// authorizedWebhoook.POST("/pix-confirm", controller.ConfirmPix)
	}

	efi.InitializePayment()
	go modules.OrderService.AnalyzePendingOrders()

	ginMode := os.Getenv("GIN_MODE")
	// frontEntEndpoint := os.Getenv("CORS_FRONTEND")
	// enableCors := os.Getenv("ENABLE_CORS")
	gin.SetMode(ginMode)

	// config := cors.DefaultConfig()

	// var corsString string = frontEntEndpoint
	// if enableCors == "false" {
	// 	corsString = "*"
	// }

	// config.AllowOrigins = []string{corsString}

	// router.Use(cors.New(config))

	router.SetTrustedProxies([]string{"127.0.0.1"})

	port := os.Getenv("GIN_PORT")
	router.Run(":" + port)
}

func Cors() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "content-type, Accept, authorization")
		c.Next()
	}
}

func Options(c *gin.Context) {

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}
