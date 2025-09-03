package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Uuid            uuid.UUID `gorm:"type:uuid;unique;not null"`
	FirstName       string    `gorm:"type:varchar(100);not null"`
	LastName        string    `gorm:"type:varchar(100);not null"`
	OIB             string    `gorm:"type:char(11);unique;not null"`
	BirthDate       time.Time `gorm:"type:date;not null"`
	Gender          string    `gorm:"type:char(1);not null"`
	MedicalRecordID uint      `gorm:"type:uint;not null"`
	MedicalRecord   MedicalRecord
	DoctorID        *uint `gorm:"type:uint;null"`
	Doctor          User
}

func (p *Patient) UpdatePatient(patient *Patient) *Patient {
	p.BirthDate = patient.BirthDate
	p.FirstName = patient.FirstName
	p.LastName = patient.LastName
	p.OIB = patient.OIB
	p.Gender = patient.Gender

	return p
}
