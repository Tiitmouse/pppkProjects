package dto

import "time"

type PatientDto struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	OIB       string    `json:"oib"`
	BirthDate time.Time `json:"birthDate"`
	Gender    string    `json:"gender"`
}

type NewPatientDto struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	OIB       string `json:"oib" binding:"required"`
	BirthDate string `json:"birthDate" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
}
