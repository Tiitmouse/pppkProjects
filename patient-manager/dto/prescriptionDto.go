package dto

import (
	"PatientManager/model"
	"time"

	"github.com/google/uuid"
)

type PrescriptionListDto struct {
	Uuid        uuid.UUID       `json:"uuid"`
	IssuedAt    time.Time       `json:"issuedAt"`
	Medications []MedicationDto `json:"medications"`
}

func (dto *PrescriptionListDto) FromModel(p *model.Prescription) *PrescriptionListDto {
	medications := make([]MedicationDto, len(p.Medications))
	for i, m := range p.Medications {
		dto := MedicationDto{}
		medications[i] = *dto.FromModel(&m)
	}

	return &PrescriptionListDto{
		Uuid:        p.Uuid,
		IssuedAt:    p.IssuedAt,
		Medications: medications,
	}
}

type CreatePrescriptionDto struct {
	IssuedAt        time.Time `json:"issuedAt" binding:"required"`
	IllnessID       uint      `json:"illnessId" binding:"required"`
	MedicationUuids []string  `json:"medicationUuids"`
}

func (dto *CreatePrescriptionDto) ToModel() *model.Prescription {
	return &model.Prescription{
		IssuedAt:  dto.IssuedAt,
		IllnessID: dto.IllnessID,
	}
}
