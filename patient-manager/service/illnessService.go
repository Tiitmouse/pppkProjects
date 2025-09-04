package service

import (
	"PatientManager/app"
	"PatientManager/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IIllnessService interface {
	Create(illness *model.Illness, recordUuid string) (*model.Illness, error)
	GetAllForRecord(recordUuid uuid.UUID) ([]model.Illness, error)
	Update(illnessUuid uuid.UUID, illnessUpdateData *model.Illness) (*model.Illness, error)
	Delete(illnessUuid uuid.UUID) error
}

type IllnessService struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewIllnessService() IIllnessService {
	var service IIllnessService
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &IllnessService{
			db:     db,
			logger: logger,
		}
	})
	return service
}

func (s *IllnessService) findMedicalRecordByUUID(recordUuid string) (*model.MedicalRecord, error) {
	var medicalRecord model.MedicalRecord
	if err := s.db.Where("uuid = ?", recordUuid).First(&medicalRecord).Error; err != nil {
		s.logger.Errorf("Error finding medical record with UUID %s: %v", recordUuid, err)
		return nil, err
	}
	return &medicalRecord, nil
}

func (s *IllnessService) Create(illness *model.Illness, recordUuid string) (*model.Illness, error) {
	illness.Uuid = uuid.New()
	medicalRecord, err := s.findMedicalRecordByUUID(recordUuid)
	if err != nil {
		return nil, err
	}
	illness.MedicalRecordID = medicalRecord.ID

	if err := s.db.Create(illness).Error; err != nil {
		s.logger.Errorf("Error creating illness: %v", err)
		return nil, err
	}
	return illness, nil
}

func (s *IllnessService) GetAllForRecord(recordUuid uuid.UUID) ([]model.Illness, error) {
	var illnesses []model.Illness
	if err := s.db.Joins("JOIN medical_records ON medical_records.id = illnesses.medical_record_id").
		Where("medical_records.uuid = ?", recordUuid).
		Order("start_date desc").
		Find(&illnesses).Error; err != nil {
		s.logger.Errorf("Error fetching illnesses for record UUID %s: %v", recordUuid, err)
		return nil, err
	}
	return illnesses, nil
}

func (s *IllnessService) findByUuid(illnessUuid uuid.UUID) (*model.Illness, error) {
	var illness model.Illness
	if err := s.db.Where("uuid = ?", illnessUuid).First(&illness).Error; err != nil {
		return nil, err
	}
	return &illness, nil
}

func (s *IllnessService) Update(illnessUuid uuid.UUID, illnessUpdateData *model.Illness) (*model.Illness, error) {
	existingIllness, err := s.findByUuid(illnessUuid)
	if err != nil {
		return nil, err
	}
	existingIllness.UpdateIllness(illnessUpdateData)
	if err := s.db.Save(existingIllness).Error; err != nil {
		s.logger.Errorf("Error saving updated illness with UUID %s: %v", illnessUuid, err)
		return nil, err
	}
	return existingIllness, nil
}

func (s *IllnessService) Delete(illnessUuid uuid.UUID) error {
	if err := s.db.Where("uuid = ?", illnessUuid).Delete(&model.Illness{}).Error; err != nil {
		s.logger.Errorf("Error deleting illness with UUID %s: %v", illnessUuid, err)
		return err
	}
	return nil
}
