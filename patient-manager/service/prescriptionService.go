package service

import (
	"PatientManager/app"
	"PatientManager/model"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IPrescriptionService interface {
	Create(prescription *model.Prescription, medicationUuids []string) (*model.Prescription, error)
	GetAllForIllness(illnessId uint) ([]model.Prescription, error)
	Delete(prescriptionUuid uuid.UUID) error
}

type PrescriptionService struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPrescriptionService() IPrescriptionService {
	var service IPrescriptionService
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &PrescriptionService{
			db:     db,
			logger: logger,
		}
	})
	return service
}

func (s *PrescriptionService) Create(prescription *model.Prescription, medicationUuids []string) (*model.Prescription, error) {
	prescription.Uuid = uuid.New()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(prescription).Error; err != nil {
			s.logger.Errorf("Error creating prescription record: %v", err)
			return err
		}

		var medications []model.Medication
		if len(medicationUuids) > 0 {
			if err := tx.Where("uuid IN ?", medicationUuids).Find(&medications).Error; err != nil {
				s.logger.Errorf("Error finding medications by UUIDs: %v", err)
				return err
			}
		}

		if len(medications) != len(medicationUuids) {
			s.logger.Error("Could not find all medications for prescription")
			return errors.New("one or more medications not found")
		}

		for i := range medications {
			medications[i].PrescriptionID = &prescription.ID
		}

		if err := tx.Save(&medications).Error; err != nil {
			s.logger.Errorf("Error updating medications with prescription ID: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.db.Preload("Medications").First(prescription, prescription.ID)
	return prescription, nil
}

func (s *PrescriptionService) GetAllForIllness(illnessId uint) ([]model.Prescription, error) {
	var prescriptions []model.Prescription
	if err := s.db.Preload("Medications").Where("illness_id = ?", illnessId).Order("issued_at desc").Find(&prescriptions).Error; err != nil {
		s.logger.Errorf("Error fetching prescriptions for illness ID %d: %v", illnessId, err)
		return nil, err
	}
	return prescriptions, nil
}

func (s *PrescriptionService) Delete(prescriptionUuid uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var prescription model.Prescription
		if err := tx.Where("uuid = ?", prescriptionUuid).First(&prescription).Error; err != nil {
			s.logger.Errorf("Error finding prescription to delete: %v", err)
			return err
		}

		if err := tx.Model(&model.Medication{}).Where("prescription_id = ?", prescription.ID).Update("prescription_id", nil).Error; err != nil {
			s.logger.Errorf("Error disassociating medications: %v", err)
			return err
		}

		if err := tx.Where("uuid = ?", prescriptionUuid).Delete(&model.Prescription{}).Error; err != nil {
			s.logger.Errorf("Error deleting prescription: %v", err)
			return err
		}

		return nil
	})
}
