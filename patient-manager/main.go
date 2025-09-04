package main

import (
	"PatientManager/app"
	"PatientManager/config"
	"PatientManager/httpServer"
	"PatientManager/repository"
	"PatientManager/service"
	"PatientManager/util/seed"

	"go.uber.org/zap"
)

// @title						PatientManager API
// @version					1.0
// @description				This is the API for the PatientManager service.
// @host						localhost:8080
// @BasePath					/api
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	app.Setup()

	// Provide logger
	app.Provide(zap.S)

	// Provide User and Login dependencies
	app.Provide(service.NewUserCrudService)
	app.Provide(service.NewLoginService)

	// Provide Patient dependencies
	app.Provide(repository.NewPatientRepository)
	app.Provide(service.NewPatientService)
	app.Provide(service.NewMedicalRecordService)
	app.Provide(service.NewChekupService)
	app.Provide(service.NewMedicationService)
	app.Provide(service.NewIllnessService)
	app.Provide(service.NewPrescriptionService)

	zap.S().Infof("Database: http://localhost:8080")

	seed.Insert()

	httpServer.Start()
}
