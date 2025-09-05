package httpServer

import (
	"PatientManager/controller"

	"github.com/gin-gonic/gin"
)

func setupHandlers(router *gin.Engine) {

	router.Static("/uploads", "./uploads")

	basePath := router.Group("/api")

	controller.NewLoginController().RegisterEndpoints(basePath)
	controller.NewPatientController().RegisterEndpoints(basePath)
	controller.NewUserController().RegisterEndpoints(basePath)
	controller.NewCheckupController().RegisterEndpoints(basePath)
	controller.NewIllnessController().RegisterEndpoints(basePath)
	controller.NewPrescriptionController().RegisterEndpoints(basePath)
	controller.NewMedicationController().RegisterEndpoints(basePath)
}
