package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Medication struct {
	gorm.Model
	Uuid           uuid.UUID `gorm:"type:uuid;unique;not null"`
	Name           string    `gorm:"type:varchar(100);not null"`
	PrescriptionID *uint     `gorm:"type:uint;null"`
	Prescription   Prescription
}

func (m *Medication) UpdateMedication(medication *Medication) *Medication {
	m.Name = medication.Name

	return m
}
