package dto

import (
	"PatientManager/model"
	"time"

	"github.com/google/uuid"
)

type PrescriptionDto struct {
	Uuid        uuid.UUID       `json:"uuid"`
	IssuedAt    time.Time       `json:"issuedAt"`
	IllnessID   uint            `json:"illnessId"`
	Medications []MedicationDto `json:"medications"`
}

func (dto *PrescriptionDto) FromModel(p *model.Prescription) *PrescriptionDto {
	medications := make([]MedicationDto, len(p.Medications))
	for i, m := range p.Medications {
		medications[i] = *(&MedicationDto{}).FromModel(&m)
	}

	return &PrescriptionDto{
		Uuid:        p.Uuid,
		IssuedAt:    p.IssuedAt,
		IllnessID:   p.IllnessID,
		Medications: medications,
	}
}

func (dto *PrescriptionDto) ToModel() (*model.Prescription, error) {
	return &model.Prescription{
		Uuid:      dto.Uuid,
		IssuedAt:  dto.IssuedAt,
		IllnessID: dto.IllnessID,
	}, nil
}
