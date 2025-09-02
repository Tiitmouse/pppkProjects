package service

import (
	"PatientManager/dto"
	"PatientManager/model"
	"PatientManager/repository"
)

type PatientService struct {
	patientRepository repository.PatientRepository
}

func NewPatientService(patientRepository repository.PatientRepository) *PatientService {
	return &PatientService{patientRepository: patientRepository}
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
	patient := model.Patient{
		FirstName: newPatient.FirstName,
		LastName:  newPatient.LastName,
		OIB:       newPatient.OIB,
		BirthDate: newPatient.BirthDate,
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
