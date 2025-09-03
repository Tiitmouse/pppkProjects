export interface PatientDto {
  id: number;
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string; 
  gender: string;
}

export interface NewPatientDto {
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string;
  gender: string;
}