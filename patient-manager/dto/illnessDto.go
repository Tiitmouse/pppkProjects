package dto

import (
	"PatientManager/model"
	"time"

	"github.com/google/uuid"
)

type IllnessDto struct {
	Uuid            uuid.UUID         `json:"uuid"`
	Name            uint              `json:"name"`
	StartDate       time.Time         `json:"startDate"`
	EndDate         *time.Time        `json:"endDate"`
	MedicalRecordID uint              `json:"medicalRecordId"`
	Prescriptions   []PrescriptionDto `json:"prescriptions"`
}

func (dto *IllnessDto) FromModel(i *model.Illness) *IllnessDto {
	prescriptions := make([]PrescriptionDto, len(i.Prescriptions))
	for j, p := range i.Prescriptions {
		prescriptions[j] = *(&PrescriptionDto{}).FromModel(&p)
	}

	return &IllnessDto{
		Uuid:            i.Uuid,
		Name:            i.Name,
		StartDate:       i.StartDate,
		EndDate:         i.EndDate,
		MedicalRecordID: i.MedicalRecordID,
		Prescriptions:   prescriptions,
	}
}

func (dto *IllnessDto) ToModel() (*model.Illness, error) {
	return &model.Illness{
		Uuid:            dto.Uuid,
		Name:            dto.Name,
		StartDate:       dto.StartDate,
		EndDate:         dto.EndDate,
		MedicalRecordID: dto.MedicalRecordID,
	}, nil
}
