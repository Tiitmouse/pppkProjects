package httpServer

import (
	"PatientManager/controller"
	"PatientManager/util/middleware"

	"github.com/gin-gonic/gin"
)

func setupHandlers(router *gin.Engine) {
	router.Use(middleware.CorsHeader())

	basePath := router.Group("/api")

	controller.NewLoginController().RegisterEndpoints(basePath)
	controller.NewPatientController().RegisterEndpoints(basePath)
	controller.NewUserController().RegisterEndpoints(basePath)
	controller.NewCheckupController().RegisterEndpoints(basePath)
	controller.NewIllnessController().RegisterEndpoints(basePath)
}
