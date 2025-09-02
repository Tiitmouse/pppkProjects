package httpServer

import (
	"PatientManager/app"
	"PatientManager/controller"
	"PatientManager/util/middleware"

	"github.com/gin-gonic/gin"
)

func patientRoutes(rg *gin.RouterGroup) {
	app.Invoke(func(patientController *controller.PatientController) {
		patient := rg.Group("/patients")
		patient.Use(middleware.Protect())
		{
			patient.GET("", patientController.GetAllPatients)
			patient.POST("", patientController.CreatePatient)
			patient.GET("/:id", patientController.GetPatientById)
			patient.PUT("/:id", patientController.UpdatePatient)
			patient.DELETE("/:id", patientController.DeletePatient)
		}
	})
}
