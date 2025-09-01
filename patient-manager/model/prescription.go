package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Prescription struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"type:uuid;unique;not null"`
	IssuedAt    time.Time `gorm:"type:date;not null"`
	IllnessID   uint      `gorm:"type:uint;not null"`
	Illness     Illness
	Medications []Medication `gorm:"foreignKey:PrescriptionID"`
}

func (p *Prescription) UpdatePrescription(prescription *Prescription) *Prescription {
	p.IssuedAt = prescription.IssuedAt

	return p
}
