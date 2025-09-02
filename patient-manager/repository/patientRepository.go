package repository

import (
	"PatientManager/model"

	"gorm.io/gorm"
)

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) FindAll() ([]model.Patient, error) {
	var patients []model.Patient
	err := r.db.Find(&patients).Error
	return patients, err
}

func (r *PatientRepository) FindById(id uint) (model.Patient, error) {
	var patient model.Patient
	err := r.db.First(&patient, id).Error
	return patient, err
}

func (r *PatientRepository) Create(patient model.Patient) (model.Patient, error) {
	err := r.db.Create(&patient).Error
	return patient, err
}

func (r *PatientRepository) Update(patient model.Patient) (model.Patient, error) {
	err := r.db.Save(&patient).Error
	return patient, err
}

func (r *PatientRepository) Delete(id uint) error {
	return r.db.Delete(&model.Patient{}, id).Error
}
