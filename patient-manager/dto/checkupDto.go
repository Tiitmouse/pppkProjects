package dto

import (
	"PatientManager/model"
	"time"

	"github.com/google/uuid"
)

// Dodan novi ImageDto
type ImageDto struct {
	Uuid string `json:"uuid"`
	Path string `json:"path"`
}

type CheckupDto struct {
	Uuid              uuid.UUID         `json:"uuid"`
	CheckupDate       time.Time         `json:"checkupDate"`
	Type              model.CheckupType `json:"type"`
	MedicalRecordUuid string            `json:"medicalRecordUuid"`
	IllnessID         *uint             `json:"illnessId,omitempty"`
	Images            []ImageDto        `json:"images"`
}

func (dto *CheckupDto) FromModel(c *model.Checkup) *CheckupDto {
	var recordUuid string
	if c.MedicalRecord.Uuid != uuid.Nil {
		recordUuid = c.MedicalRecord.Uuid.String()
	}

	// Kreiranje ImageDto-a
	imageDtos := make([]ImageDto, len(c.Images))
	for i, image := range c.Images {
		imageDtos[i] = ImageDto{
			Uuid: image.Uuid.String(),
			Path: image.Path,
		}
	}

	return &CheckupDto{
		Uuid:              c.Uuid,
		CheckupDate:       c.CheckupDate,
		Type:              c.Type,
		MedicalRecordUuid: recordUuid,
		IllnessID:         c.IllnessID,
		Images:            imageDtos, // Dodano
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
	Images            []string          `json:"images"`
}

func (dto *CreateCheckupDto) ToModel() (*model.Checkup, error) {
	return &model.Checkup{
		CheckupDate: dto.CheckupDate,
		Type:        dto.Type,
		IllnessID:   dto.IllnessID,
	}, nil
}
