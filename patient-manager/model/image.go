package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Uuid      uuid.UUID `gorm:"type:uuid;unique;not null"`
	Path      string    `gorm:"type:varchar(255);not null"`
	CheckupID uint
	Checkup   Checkup
}
