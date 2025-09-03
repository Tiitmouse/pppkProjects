package dto

import (
	"PatientManager/model"
	"PatientManager/util/cerror"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NewUserDto struct {
	Uuid      string `json:"uuid"`
	FirstName string `json:"firstName" binding:"required,min=2,max=100"`
	LastName  string `json:"lastName" binding:"required,min=2,max=100"`
	OIB       string `json:"oib" binding:"required,len=11"`
	BirthDate string `json:"birthDate" binding:"required,datetime=2006-01-02"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password"`
	Role      string `json:"role" binding:"required,oneof=doctor patient superadmin"`
}

// ToModel create a model from a dto
func (dto *NewUserDto) ToModel() (*model.User, error) {
	role, err := model.StoUserRole(dto.Role)
	if err != nil {
		zap.S().Error("Failed to parse role = %+v, err = %+v", dto.Role, err)
		return nil, cerror.ErrUnknownRole
	}
	if dto.Uuid != "" {
		_, err := uuid.Parse(dto.Uuid)
		if err != nil {
			zap.S().Errorf("Failed to parse uuid = %s, err = %+v", dto.Uuid, err)
			return nil, cerror.ErrBadUuid
		}
	}

	return &model.User{
		Uuid:      uuid.New(),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Role:      role,
		OIB:       dto.OIB,
	}, nil
}

// FromModel returns a dto from model struct
func (dto *NewUserDto) FromModel(m *model.User) *NewUserDto {
	dto = &NewUserDto{
		Uuid:      m.Uuid.String(),
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		Role:      fmt.Sprint(m.Role),
	}
	return dto
}
