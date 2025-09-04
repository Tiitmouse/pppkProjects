package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckupType string

const (
	GeneralPractitioner CheckupType = "GP"
	BloodTest           CheckupType = "KRV"
	XRayScan            CheckupType = "X-RAY"
	CTScan              CheckupType = "CT"
	MRIScan             CheckupType = "MR"
	Ultrasound          CheckupType = "ULTRA"
	Electrocardiogram   CheckupType = "EKG"
	Echocardiogram      CheckupType = "ECHO"
	EyeExam             CheckupType = "EYE"
	DermatologyExam     CheckupType = "DERM"
	DentalExam          CheckupType = "DENTA"
	Mammography         CheckupType = "MAMMO"
	NeurologyExam       CheckupType = "NEURO"
)

type Checkup struct {
	gorm.Model
	Uuid            uuid.UUID   `gorm:"type:uuid;unique;not null"`
	CheckupDate     time.Time   `gorm:"not null"`
	Type            CheckupType `gorm:"type:varchar(10);not null"`
	MedicalRecordID uint        `gorm:"type:uint;not null"`
	MedicalRecord   MedicalRecord
	IllnessID       *uint `gorm:"type:uint;null"`
	Illness         Illness
	Images          []Image `gorm:"foreignKey:CheckupID"`
}

func (c *Checkup) UpdateCheckup(checkup *Checkup) *Checkup {
	c.CheckupDate = checkup.CheckupDate
	c.Type = checkup.Type
	c.IllnessID = checkup.IllnessID

	return c
}
