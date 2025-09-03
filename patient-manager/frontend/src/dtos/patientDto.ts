import type { Patient } from '@/models/patient';
import { formatDate } from '@/utils/formatDate';
import type { UserRole } from '@/enums/userRole';


export interface PatientDto {
  id: number;
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string; 
  gender: string;
}

export interface NewPatientDto {
  uuid: string
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string;
  gender: string;
  role: UserRole.Patient
}

export function createNewPatientDto(
    patient: Patient
): NewPatientDto {
    return {
        uuid: patient.uuid,
        firstName: patient.firstName,
        lastName: patient.lastName,
        oib: patient.oib,
        gender: patient.gender,
        birthDate: formatDate(patient.birthDate),
        role: patient.role
    };
}