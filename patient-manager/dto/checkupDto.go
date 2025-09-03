package dto

import (
	"PatientManager/model"
	"time"

	"github.com/google/uuid"
)

type CheckupDto struct {
	Uuid            uuid.UUID         `json:"uuid"`
	PatientID       uint              `json:"patientId"`
	CheckupDate     time.Time         `json:"checkupDate"`
	Type            model.CheckupType `json:"type"`
	MedicalRecordID uint              `json:"medicalRecordId"`
	IllnessID       *uint             `json:"illnessId"`
}

func (dto *CheckupDto) FromModel(c *model.Checkup) *CheckupDto {
	return &CheckupDto{
		Uuid:            c.Uuid,
		CheckupDate:     c.CheckupDate,
		Type:            c.Type,
		MedicalRecordID: c.MedicalRecordID,
		IllnessID:       c.IllnessID,
	}
}

func (dto *CheckupDto) ToModel() (*model.Checkup, error) {
	return &model.Checkup{
		Uuid:            dto.Uuid,
		CheckupDate:     dto.CheckupDate,
		Type:            dto.Type,
		MedicalRecordID: dto.MedicalRecordID,
		IllnessID:       dto.IllnessID,
	}, nil
}
