package seed

import (
	"PatientManager/app"
	"PatientManager/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func seedMedications() error {
	var err error
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		medications := []string{
			"Aspirin", "Paracetamol", "Ibuprofen", "Amoxicillin", "Lisinopril",
			"Atorvastatin", "Metformin", "Simvastatin", "Omeprazole", "Amlodipine",
			"Metoprolol", "Acetaminophen", "Hydrochlorothiazide", "Sertraline",
			"Citalopram", "Zolpidem", "Furosemide", "Alprazolam", "Escitalopram",
		}

		var count int64
		db.Model(&model.Medication{}).Count(&count)
		if count > 0 {
			logger.Infoln("Medications already seeded. Skipping.")
			return
		}

		logger.Infoln("Seeding medications...")
		for _, name := range medications {
			medication := model.Medication{
				Uuid: uuid.New(),
				Name: name,
			}
			if creationErr := db.Create(&medication).Error; creationErr != nil {
				logger.Errorf("Failed to seed medication %s: %v", name, creationErr)
				err = creationErr
				return
			}
		}

		logger.Info("Medication seeding completed successfully.")
	})

	return err
}
