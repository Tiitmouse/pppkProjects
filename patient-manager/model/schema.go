package model

// GetAllModels returns an array of all registered models
func GetAllModels() []any {
	return []any{
		&User{},
		&Patinet{},
		&MedicalRecord{},
		&Checkup{},
		&Prescription{},
		&Medication{},
		&Illness{},
	}
}
