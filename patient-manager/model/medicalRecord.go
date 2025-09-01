package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicalRecord struct {
	gorm.Model
	Uuid      uuid.UUID `gorm:"type:uuid;unique;not null"`
	PatientID uint      `gorm:"type:uint;not null"`
	DoctorID  uint      `gorm:"type:uint;not null"`
	Checkups  []Checkup `gorm:"foreignKey:MedicalRecordID"`
	Illnesses []Illness `gorm:"foreignKey:MedicalRecordID"`
}

func (mr *MedicalRecord) UpdateMedicalRecord(medicalRecord *MedicalRecord) *MedicalRecord {
	mr.DoctorID = medicalRecord.DoctorID

	return mr
}
