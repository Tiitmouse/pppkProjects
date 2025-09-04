package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Illness struct {
	gorm.Model
	Uuid            uuid.UUID  `gorm:"type:uuid;unique;not null"`
	Name            string     `gorm:"type:varchar(100);not null"`
	StartDate       time.Time  `gorm:"type:date;not null"`
	EndDate         *time.Time `gorm:"type:date;null"`
	MedicalRecordID uint       `gorm:"type:uint;not null"`
	MedicalRecord   MedicalRecord
	Prescriptions   []Prescription `gorm:"foreignKey:IllnessID"`
}

func (i *Illness) UpdateIllness(illness *Illness) *Illness {
	i.Name = illness.Name
	i.StartDate = illness.StartDate
	i.EndDate = illness.EndDate

	return i
}
