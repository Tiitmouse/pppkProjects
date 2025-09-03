package service

import (
	"PatientManager/app"
	"PatientManager/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IMedicalRecordService interface {
	Create(record *model.MedicalRecord) (*model.MedicalRecord, error)
	Read(patientOib string) (*model.MedicalRecord, error)
	Update(recordUuid uuid.UUID, recordUpdateData *model.MedicalRecord) (*model.MedicalRecord, error)
	Delete(recordUuid uuid.UUID) error
}

type MedicalRecordService struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewMedicalRecordService() IMedicalRecordService {
	var service IMedicalRecordService
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &MedicalRecordService{
			db:     db,
			logger: logger,
		}
	})

	return service
}

func (s *MedicalRecordService) Create(record *model.MedicalRecord) (*model.MedicalRecord, error) {
	record.Uuid = uuid.New()
	s.logger.Infof("Creating medical record for patient ID: %d", record.PatientID)

	rez := s.db.Create(record)
	if rez.Error != nil {
		s.logger.Errorf("Error creating medical record: %v", rez.Error)
		return nil, rez.Error
	}

	s.logger.Infof("Successfully created medical record with UUID: %s", record.Uuid)
	return record, nil
}

func (s *MedicalRecordService) Read(patientOib string) (*model.MedicalRecord, error) {
	var patient model.Patient
	if err := s.db.Where("oib = ?", patientOib).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warnf("Patient with OIB %s not found", patientOib)
		} else {
			s.logger.Errorf("Error finding patient with OIB %s: %v", patientOib, err)
		}
		return nil, err
	}

	var medicalRecord model.MedicalRecord
	rez := s.db.
		Preload("Checkups").
		Preload("Illnesses").
		Where("patient_id = ?", patient.ID).
		First(&medicalRecord)

	if rez.Error != nil {
		if rez.Error == gorm.ErrRecordNotFound {
			s.logger.Warnf("Medical record for patient with OIB %s not found", patientOib)
		} else {
			s.logger.Errorf("Error finding medical record for patient with OIB %s: %v", patientOib, rez.Error)
		}
		return nil, rez.Error
	}

	return &medicalRecord, nil
}

func (s *MedicalRecordService) findByUuid(recordUuid uuid.UUID) (*model.MedicalRecord, error) {
	var record model.MedicalRecord
	rez := s.db.
		Where("uuid = ?", recordUuid).
		First(&record)

	if rez.Error != nil {
		if rez.Error == gorm.ErrRecordNotFound {
			s.logger.Warnf("Medical record with UUID %s not found", recordUuid)
		} else {
			s.logger.Errorf("Error finding medical record with UUID %s: %v", recordUuid, rez.Error)
		}
		return nil, rez.Error
	}
	return &record, nil
}

func (s *MedicalRecordService) Update(recordUuid uuid.UUID, recordUpdateData *model.MedicalRecord) (*model.MedicalRecord, error) {
	existingRecord, err := s.findByUuid(recordUuid)
	if err != nil {
		return nil, err
	}

	s.logger.Debugf("Updating medical record with UUID: %s", recordUuid)

	existingRecord.UpdateMedicalRecord(recordUpdateData)

	rez := s.db.Save(existingRecord)
	if rez.Error != nil {
		s.logger.Errorf("Error saving updated medical record with UUID %s: %v", recordUuid, rez.Error)
		return nil, rez.Error
	}

	s.logger.Infof("Successfully updated medical record with UUID: %s", recordUuid)
	return existingRecord, nil
}

func (s *MedicalRecordService) Delete(recordUuid uuid.UUID) error {
	s.logger.Infof("Attempting to delete medical record with UUID: %s", recordUuid)

	rez := s.db.Where("uuid = ?", recordUuid).Delete(&model.MedicalRecord{})

	if rez.Error != nil {
		s.logger.Errorf("Error deleting medical record with UUID %s: %v", recordUuid, rez.Error)
		return rez.Error
	}

	if rez.RowsAffected == 0 {
		s.logger.Warnf("No medical record found with UUID %s to delete", recordUuid)
		return gorm.ErrRecordNotFound
	}

	s.logger.Infof("Successfully soft-deleted medical record with UUID: %s", recordUuid)
	return nil
}
