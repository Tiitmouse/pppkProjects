package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Illness struct {
	gorm.Model
	Uuid            uuid.UUID  `gorm:"type:uuid;unique;not null"`
	Name            uint       `gorm:"type:varchar(100);not null"`
	StartDate       time.Time  `gorm:"type:date;not null"`
	EndDate         *time.Time `gorm:"type:date;null"`
	MedicalRecordID uint       `gorm:"type:uint;not null"`
	MedicalRecord   MedicalRecord
	Prescriptions   []Prescription `gorm:"foreignKey:IllnessID"`
}

func (i *Illness) UpdateIllness(illnes *Illness) *Illness {
	i.Name = illnes.Name
	i.StartDate = illnes.StartDate
	i.EndDate = illnes.EndDate

	return i
}
