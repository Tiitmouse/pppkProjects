package main

import (
	"PatientManager/app"
	"PatientManager/config"
	"PatientManager/controller"
	"PatientManager/httpServer"
	"PatientManager/repository"
	"PatientManager/service"
	"PatientManager/util/seed"

	"go.uber.org/zap"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	app.Setup()

	// Provide logger
	app.Provide(zap.S)

	app.Provide(service.NewLoginService)
	app.Provide(service.NewUserCrudService)
	app.Provide(controller.NewLoginController)
	app.Provide(controller.NewUserController)

	app.Provide(repository.NewPatientRepository)
	app.Provide(service.NewPatientService)
	app.Provide(controller.NewPatientController)

	zap.S().Infof("Database: http://localhost:8080")

	seed.Insert()

	httpServer.Start()
}
