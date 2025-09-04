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
	patientRepository    repository.PatientRepository
	medicalRecordService IMedicalRecordService
}

type IPatientService interface {
	GetAllPatients() ([]dto.PatientDto, error)
	GetPatientById(id uint) (dto.PatientDto, error)
	CreatePatient(newPatient dto.NewPatientDto) (dto.PatientDto, error)
	UpdatePatient(id uint, patientDto dto.UpdatePatientDto) (dto.PatientDto, error)
	DeletePatient(id uint) error
}

func NewPatientService() IPatientService {
	var service *PatientService
	app.Invoke(func(repo repository.PatientRepository, mrservice IMedicalRecordService) {
		service = &PatientService{
			patientRepository:    repo,
			medicalRecordService: mrservice,
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
		patientDtos = append(patientDtos, dto.FromModel(&p))
	}
	return patientDtos, nil
}

func (s *PatientService) GetPatientById(id uint) (dto.PatientDto, error) {
	patient, err := s.patientRepository.FindByIdWithDoctor(id)
	if err != nil {
		return dto.PatientDto{}, err
	}
	return dto.FromModel(&patient), nil
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
		DoctorID:  newPatient.DoctorID,
	}

	createdPatient, err := s.patientRepository.Create(patient)
	if err != nil {
		return dto.PatientDto{}, err
	}

	medicalRecord := model.MedicalRecord{
		PatientID: createdPatient.ID,
	}

	createdmr, err := s.medicalRecordService.Create(&medicalRecord)
	if err != nil {
		return dto.PatientDto{}, err
	}

	createdPatient.MedicalRecordID = createdmr.ID
	createdPatient.MedicalRecord = *createdmr

	s.patientRepository.Update(createdPatient)

	return dto.FromModel(&createdPatient), nil
}

func (s *PatientService) UpdatePatient(id uint, patientDto dto.UpdatePatientDto) (dto.PatientDto, error) {
	patient, err := s.patientRepository.FindById(id)
	if err != nil {
		return dto.PatientDto{}, err
	}

	patient.FirstName = patientDto.FirstName
	patient.LastName = patientDto.LastName
	patient.OIB = patientDto.OIB

	bod, err := time.Parse(time.RFC3339, patientDto.BirthDate)
	if err != nil {
		zap.S().Errorf("Failed to parse BirthDate = %s, err = %+v", patientDto.BirthDate, err)
		return dto.PatientDto{}, cerror.ErrBadDateFormat
	}

	patient.BirthDate = bod
	patient.Gender = patientDto.Gender
	patient.DoctorID = patientDto.DoctorID

	updatedPatient, err := s.patientRepository.Update(patient)
	if err != nil {
		return dto.PatientDto{}, err
	}

	return dto.FromModel(&updatedPatient), nil
}

func (s *PatientService) DeletePatient(id uint) error {
	return s.patientRepository.Delete(id)
}
