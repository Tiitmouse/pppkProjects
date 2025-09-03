package dto

import (
	"PatientManager/model"

	"github.com/google/uuid"
)

type MedicalRecordDto struct {
	Uuid      uuid.UUID    `json:"uuid"`
	PatientID uint         `json:"patientId"`
	DoctorID  uint         `json:"doctorId"`
	Checkups  []CheckupDto `json:"checkups"`
	Illnesses []IllnessDto `json:"illnesses"`
}

func (dto *MedicalRecordDto) FromModel(mr *model.MedicalRecord) *MedicalRecordDto {
	checkups := make([]CheckupDto, len(mr.Checkups))
	for i, checkup := range mr.Checkups {
		checkups[i] = *(&CheckupDto{}).FromModel(&checkup)
	}

	illnesses := make([]IllnessDto, len(mr.Illnesses))
	for i, illness := range mr.Illnesses {
		illnesses[i] = *(&IllnessDto{}).FromModel(&illness)
	}

	return &MedicalRecordDto{
		Uuid:      mr.Uuid,
		PatientID: mr.PatientID,
		DoctorID:  mr.DoctorID,
		Checkups:  checkups,
		Illnesses: illnesses,
	}
}

func (dto *MedicalRecordDto) ToModel() (*model.MedicalRecord, error) {
	return &model.MedicalRecord{
		Uuid:      dto.Uuid,
		PatientID: dto.PatientID,
		DoctorID:  dto.DoctorID,
	}, nil
}

type NewMedicalRecordDto struct {
	PatientID uint `json:"patientId"`
	DoctorID  uint `json:"doctorId"`
}

func (dto *NewMedicalRecordDto) ToModel() (*model.MedicalRecord, error) {
	return &model.MedicalRecord{
		Uuid:      uuid.New(),
		PatientID: dto.PatientID,
		DoctorID:  dto.DoctorID,
	}, nil
}
