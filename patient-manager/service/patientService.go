package service

import (
	"PatientManager/app"
	"PatientManager/dto"
	"PatientManager/model"
	"PatientManager/repository"
	"PatientManager/util/cerror"
	"PatientManager/util/format"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PatientService struct {
	patientRepository repository.PatientRepository
}

type IPatientService interface {
	GetAllPatients() ([]dto.PatientDto, error)
	GetPatientById(id uint) (dto.PatientDto, error)
	CreatePatient(newPatient dto.NewPatientDto) (dto.PatientDto, error)
	UpdatePatient(id uint, patientDto dto.PatientDto) (dto.PatientDto, error)
	DeletePatient(id uint) error
}

func NewPatientService() IPatientService {
	var service *PatientService
	app.Invoke(func(repo repository.PatientRepository) {
		service = &PatientService{
			patientRepository: repo,
		}
	})
	return service
}

func (s *PatientService) GetAllPatients() ([]dto.PatientDto, error) {
	patients, err := s.patientRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var patientDtos []dto.PatientDto
	for _, p := range patients {
		patientDtos = append(patientDtos, dto.PatientDto{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			OIB:       p.OIB,
			BirthDate: p.BirthDate,
			Gender:    p.Gender,
		})
	}
	return patientDtos, nil
}

func (s *PatientService) GetPatientById(id uint) (dto.PatientDto, error) {
	patient, err := s.patientRepository.FindById(id)
	if err != nil {
		return dto.PatientDto{}, err
	}
	return dto.PatientDto{
		ID:        patient.ID,
		FirstName: patient.FirstName,
		LastName:  patient.LastName,
		OIB:       patient.OIB,
		BirthDate: patient.BirthDate,
		Gender:    patient.Gender,
	}, nil
}

func (s *PatientService) CreatePatient(newPatient dto.NewPatientDto) (dto.PatientDto, error) {
	bod, err := time.Parse(format.DateFormat, newPatient.BirthDate)
	if err != nil {
		zap.S().Errorf("Failed to parse BirthDate = %s, err = %+v", newPatient.BirthDate, err)
		return dto.PatientDto{}, cerror.ErrBadDateFormat
	}

	patient := model.Patient{
		Uuid:      uuid.New(),
		FirstName: newPatient.FirstName,
		LastName:  newPatient.LastName,
		OIB:       newPatient.OIB,
		BirthDate: bod,
		Gender:    newPatient.Gender,
	}

	createdPatient, err := s.patientRepository.Create(patient)
	if err != nil {
		return dto.PatientDto{}, err
	}

	return dto.PatientDto{
		ID:        createdPatient.ID,
		FirstName: createdPatient.FirstName,
		LastName:  createdPatient.LastName,
		OIB:       createdPatient.OIB,
		BirthDate: createdPatient.BirthDate,
		Gender:    createdPatient.Gender,
	}, nil
}

func (s *PatientService) UpdatePatient(id uint, patientDto dto.PatientDto) (dto.PatientDto, error) {
	patient, err := s.patientRepository.FindById(id)
	if err != nil {
		return dto.PatientDto{}, err
	}

	patient.FirstName = patientDto.FirstName
	patient.LastName = patientDto.LastName
	patient.OIB = patientDto.OIB
	patient.BirthDate = patientDto.BirthDate
	patient.Gender = patientDto.Gender

	updatedPatient, err := s.patientRepository.Update(patient)
	if err != nil {
		return dto.PatientDto{}, err
	}

	return dto.PatientDto{
		ID:        updatedPatient.ID,
		FirstName: updatedPatient.FirstName,
		LastName:  updatedPatient.LastName,
		OIB:       updatedPatient.OIB,
		BirthDate: updatedPatient.BirthDate,
		Gender:    updatedPatient.Gender,
	}, nil
}

func (s *PatientService) DeletePatient(id uint) error {
	return s.patientRepository.Delete(id)
}
