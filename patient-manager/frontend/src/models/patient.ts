import type { UserRole } from '@/enums/userRole';

export interface Patient {
  uuid: string;
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string;
  gender: string;
  role: UserRole.Patient;
  medicalRecordUuid: string
}