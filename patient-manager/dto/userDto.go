package dto

import (
	"PatientManager/model"
	"PatientManager/util/cerror"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserDto struct {
	Uuid      string `json:"uuid"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	OIB       string `json:"oib"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

func (dto *UserDto) ToModel() (*model.User, error) {
	uuid, err := uuid.Parse(dto.Uuid)
	if err != nil {
		zap.S().Error("Failed to parse uuid = %s, err = %+v", dto.Uuid, err)
		return nil, cerror.ErrBadUuid
	}

	role, err := model.StoUserRole(dto.Role)
	if err != nil {
		zap.S().Errorf("Failed to parse role = %+v, err = %+v", dto.Role, err)
		return nil, cerror.ErrUnknownRole
	}

	return &model.User{
		Uuid:      uuid,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Role:      role,
	}, nil
}

// FromModel returns a dto from model struct
func (dto UserDto) FromModel(m *model.User) UserDto {
	dto = UserDto{
		Uuid:      m.Uuid.String(),
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		Role:      fmt.Sprint(m.Role),
	}
	return dto
}
