package model

import (
	"PatientManager/util/cerror"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleDoctor     UserRole = "doctor"
	RolePatient    UserRole = "patient"
	RoleSuperAdmin UserRole = "superadmin"
)

func StoUserRole(text string) (UserRole, error) {
	switch text {
	case fmt.Sprint(RoleDoctor):
		return RoleDoctor, nil

	case fmt.Sprint(RolePatient):
		return RolePatient, nil

	case fmt.Sprint(RoleSuperAdmin):
		return RoleSuperAdmin, nil

	default:
		return "", cerror.ErrUnknownRole
	}
}

type User struct {
	gorm.Model
	Uuid         uuid.UUID `gorm:"type:uuid;unique;not null"`
	FirstName    string    `gorm:"type:varchar(100);not null"`
	LastName     string    `gorm:"type:varchar(100);not null"`
	OIB          string    `gorm:"type:char(11);not null"`
	Email        string    `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         UserRole  `gorm:"type:varchar(20);not null"`
	Patients     []Patient `gorm:"foreignKey:DoctorID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	validRoles := map[UserRole]bool{
		RoleDoctor: true, RolePatient: true, RoleSuperAdmin: true,
	}

	if _, ok := validRoles[u.Role]; !ok {
		return errors.New("invalid user role")
	}
	return nil
}

func (u *User) Update(user *User) *User {
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.Role = user.Role

	return u
}
