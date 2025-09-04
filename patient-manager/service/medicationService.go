package service

import (
	"PatientManager/app"
	"PatientManager/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IMedicationService interface {
	GetAll() ([]model.Medication, error)
}

type MedicationService struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewMedicationService() IMedicationService {
	var service IMedicationService
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &MedicationService{
			db:     db,
			logger: logger,
		}
	})
	return service
}

func (s *MedicationService) GetAll() ([]model.Medication, error) {
	var medications []model.Medication
	if err := s.db.Find(&medications).Error; err != nil {
		s.logger.Errorf("Error fetching all medications: %v", err)
		return nil, err
	}
	return medications, nil
}
