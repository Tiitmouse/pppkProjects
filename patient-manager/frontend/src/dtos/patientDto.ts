import type { Patient } from '@/models/patient';
import { formatDate } from '@/utils/formatDate';
import type { DoctorDto } from './userDto';

export interface PatientDto {
  id: number;
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string; 
  gender: string;
  medicalRecordUuid: string;
  doctor?: DoctorDto;
}

export interface NewPatientDto {
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string;
  gender: string;
  doctorId?: number;
}

export function createNewPatientDto(
    patient: Patient
): NewPatientDto {
    return {
        firstName: patient.firstName,
        lastName: patient.lastName,
        oib: patient.oib,
        gender: patient.gender,
        birthDate: formatDate(patient.birthDate),
    };
}