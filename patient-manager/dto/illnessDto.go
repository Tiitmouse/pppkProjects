package dto

import (
	"PatientManager/model"
	"time"
)

type CreateIllnessDto struct {
	Name              string     `json:"name" binding:"required"`
	StartDate         time.Time  `json:"startDate" binding:"required"`
	EndDate           *time.Time `json:"endDate"`
	MedicalRecordUuid string     `json:"medicalRecordUuid" binding:"required"`
}

func (dto *CreateIllnessDto) ToModel() *model.Illness {
	return &model.Illness{
		Name:      dto.Name,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
	}
}

type UpdateIllnessDto struct {
	Name      string     `json:"name" binding:"required"`
	StartDate time.Time  `json:"startDate" binding:"required"`
	EndDate   *time.Time `json:"endDate"`
}

func (dto *UpdateIllnessDto) ToModel() *model.Illness {
	return &model.Illness{
		Name:      dto.Name,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
	}
}

type IllnessListDto struct {
	Uuid      string     `json:"uuid"`
	Name      string     `json:"name"`
	StartDate time.Time  `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
}

func (dto *IllnessListDto) FromModel(i *model.Illness) *IllnessListDto {
	return &IllnessListDto{
		Uuid:      i.Uuid.String(),
		Name:      i.Name,
		StartDate: i.StartDate,
		EndDate:   i.EndDate,
	}
}
