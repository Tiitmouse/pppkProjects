package httpServer

import (
	"PatientManager/app"
	"PatientManager/controller"
	_ "PatientManager/docs"
	"PatientManager/util/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupHandlers(router *gin.Engine) {
	router.Use(middleware.CorsHeader())

	basePath := router.Group("/api")
	basePath.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Invoke(func(
		loginController *controller.LoginController,
		userController *controller.UserController,
		patientController *controller.PatientController,
	) {
		loginController.RegisterEndpoints(basePath)
		userController.RegisterEndpoints(basePath)
		patientController.RegisterEndpoints(basePath)
	})
}
