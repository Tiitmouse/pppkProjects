package httpServer

import (
	"PatientManager/controller"
	"PatientManager/docs"
	"PatientManager/util/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupHandlers(router *gin.Engine) {
	router.Use(middleware.CorsHeader())
	api := router.Group("/api")

	// register swagger
	docs.SwaggerInfo.BasePath = "/api"

	swagger := ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8090/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(2))
	router.GET("/swagger/*any", swagger)

	// NOTE: register controllers here
	controller.NewLoginController().RegisterEndpoints(api)
	controller.NewUserController().RegisterEndpoints(api)
}
