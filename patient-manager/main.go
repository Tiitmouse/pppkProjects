package main

import (
	"PatientManager/app"
	"PatientManager/config"
	"PatientManager/httpServer"
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

	// Provided logger
	app.Provide(zap.S)

	app.Provide(service.NewLoginService)
	app.Provide(service.NewUserCrudService)

	zap.S().Infof("Database: http://localhost:8080")

	seed.Insert()

	httpServer.Start()
}
