package dto

import (
	"PatientManager/model"
	"time"

	"github.com/google/uuid"
)

type CheckupDto struct {
	Uuid              uuid.UUID         `json:"uuid"`
	CheckupDate       time.Time         `json:"checkupDate"`
	Type              model.CheckupType `json:"type"`
	MedicalRecordUuid string            `json:"medicalRecordUuid"`
	IllnessID         *uint             `json:"illnessId,omitempty"`
}

func (dto *CheckupDto) FromModel(c *model.Checkup) *CheckupDto {
	var recordUuid string
	if c.MedicalRecord.Uuid != uuid.Nil {
		recordUuid = c.MedicalRecord.Uuid.String()
	}

	return &CheckupDto{
		Uuid:              c.Uuid,
		CheckupDate:       c.CheckupDate,
		Type:              c.Type,
		MedicalRecordUuid: recordUuid,
		IllnessID:         c.IllnessID,
	}
}

func (dto *CheckupDto) ToModel() (*model.Checkup, error) {
	return &model.Checkup{
		Uuid:        dto.Uuid,
		CheckupDate: dto.CheckupDate,
		Type:        dto.Type,
		IllnessID:   dto.IllnessID,
	}, nil
}

type CreateCheckupDto struct {
	CheckupDate       time.Time         `json:"checkupDate" binding:"required"`
	Type              model.CheckupType `json:"type" binding:"required"`
	MedicalRecordUuid string            `json:"medicalRecordUuid" binding:"required"`
	IllnessID         *uint             `json:"illnessId"`
}

func (dto *CreateCheckupDto) ToModel() (*model.Checkup, error) {
	return &model.Checkup{
		CheckupDate: dto.CheckupDate,
		Type:        dto.Type,
		IllnessID:   dto.IllnessID,
	}, nil
}
