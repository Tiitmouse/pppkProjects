package model

func GetAllModels() []any {
	return []any{
		&User{},
		&Patient{},
		&MedicalRecord{},
		&Checkup{},
		&Prescription{},
		&Medication{},
		&Illness{},
		&Image{},
	}
}
