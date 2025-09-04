package dto

import (
	"PatientManager/model"
	"time"
)

type PatientDto struct {
	ID                uint       `json:"id"`
	FirstName         string     `json:"firstName"`
	LastName          string     `json:"lastName"`
	OIB               string     `json:"oib"`
	BirthDate         time.Time  `json:"birthDate"`
	Gender            string     `json:"gender"`
	MedicalRecordUuid string     `json:"medicalRecordUuid"`
	Doctor            *DoctorDto `json:"doctor,omitempty"`
}

type NewPatientDto struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	OIB       string `json:"oib" binding:"required"`
	BirthDate string `json:"birthDate" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
	DoctorID  *uint  `json:"doctorId,omitempty"`
}

type UpdatePatientDto struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	OIB       string `json:"oib"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
	DoctorID  *uint  `json:"doctorId"`
}

func FromModel(p *model.Patient) PatientDto {
	var doctorDto *DoctorDto
	if p.DoctorID != nil {
		doctorDto = &DoctorDto{
			FirstName: p.Doctor.FirstName,
			LastName:  p.Doctor.LastName,
		}
	}

	return PatientDto{
		ID:                p.ID,
		FirstName:         p.FirstName,
		LastName:          p.LastName,
		OIB:               p.OIB,
		BirthDate:         p.BirthDate,
		Gender:            p.Gender,
		MedicalRecordUuid: p.MedicalRecord.Uuid.String(),
		Doctor:            doctorDto,
	}
}
