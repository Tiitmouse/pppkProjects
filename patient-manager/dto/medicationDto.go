package dto

import (
	"PatientManager/model"

	"github.com/google/uuid"
)

type MedicationDto struct {
	Uuid           uuid.UUID `json:"uuid"`
	Name           string    `json:"name"`
	PrescriptionID uint      `json:"prescriptionId"`
}

func (dto *MedicationDto) FromModel(m *model.Medication) *MedicationDto {
	return &MedicationDto{
		Uuid:           m.Uuid,
		Name:           m.Name,
		PrescriptionID: m.PrescriptionID,
	}
}

func (dto *MedicationDto) ToModel() (*model.Medication, error) {
	return &model.Medication{
		Uuid:           dto.Uuid,
		Name:           dto.Name,
		PrescriptionID: dto.PrescriptionID,
	}, nil
}
