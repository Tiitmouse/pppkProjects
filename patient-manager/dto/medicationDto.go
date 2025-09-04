package dto

import (
	"PatientManager/model"

	"github.com/google/uuid"
)

// MedicationDto is used when viewing a medication within a prescription.
type MedicationDto struct {
	Uuid           uuid.UUID `json:"uuid"`
	Name           string    `json:"name"`
	PrescriptionID uint      `json:"prescriptionId"`
}

func (dto *MedicationDto) FromModel(m *model.Medication) *MedicationDto {
	var prescriptionID uint
	if m.PrescriptionID != nil {
		prescriptionID = *m.PrescriptionID
	}

	return &MedicationDto{
		Uuid:           m.Uuid,
		Name:           m.Name,
		PrescriptionID: prescriptionID,
	}
}

func (dto *MedicationDto) ToModel() (*model.Medication, error) {
	return &model.Medication{
		Uuid:           dto.Uuid,
		Name:           dto.Name,
		PrescriptionID: &dto.PrescriptionID,
	}, nil
}

// MedicationListDto is used for the general list of all available medications.
type MedicationListDto struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func (dto *MedicationListDto) FromModel(m *model.Medication) *MedicationListDto {
	return &MedicationListDto{
		Uuid: m.Uuid.String(),
		Name: m.Name,
	}
}
