package httpServer

import (
	"PatientManager/app"
	"PatientManager/controller"
	"PatientManager/util/middleware"

	"github.com/gin-gonic/gin"
)

func patientRoutes(rg *gin.RouterGroup) {
	app.Invoke(func(patientController *controller.PatientController, authMiddleware *middleware.AuthMiddleware) {
		patient := rg.Group("/patients")
		patient.Use(authMiddleware.Handler())
		{
			patient.GET("", patientController.GetAllPatients)
			patient.POST("", patientController.CreatePatient)
			patient.GET("/:id", patientController.GetPatientById)
			patient.PUT("/:id", patientController.UpdatePatient)
			patient.DELETE("/:id", patientController.DeletePatient)
		}
	})
}
