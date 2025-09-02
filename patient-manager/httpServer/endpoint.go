package httpServer

import (
	"PatientManager/app"
	"PatientManager/controller"
	_ "PatientManager/docs"
	"PatientManager/util/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupHandlers(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	basePath := router.Group("/api")

	// Swagger endpoint
	basePath.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Invoke a function to get controllers and set up routes
	app.Invoke(func(
		loginController *controller.LoginController,
		userController *controller.UserController,
		authMiddleware *middleware.AuthMiddleware,
	) {
		// Auth endpoints
		auth := basePath.Group("/auth")
		{
			auth.POST("/login", loginController.Login)
			auth.POST("/refresh", loginController.Refresh)
		}

		// User endpoints
		user := basePath.Group("/user")
		{
			user.POST("", userController.CreateUser)
			user.Use(authMiddleware.Handler())
			user.GET("", userController.GetAllUsers)
			user.GET("/:uuid", userController.GetUserByUuid)
			user.PUT("/:uuid", userController.UpdateUser)
			user.DELETE("/:uuid", userController.DeleteUser)
		}

		// Call the function to set up patient routes
		patientRoutes(basePath)
	})
}
