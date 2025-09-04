import { defineStore } from 'pinia'
import type { DoctorDto } from '@/dtos/userDto';

export interface Patient {
  id: number;
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string;
  gender: string;
  medicalRecordUuid: string;
  doctor?: DoctorDto;
}

export const usePatientStore = defineStore('patient', {
  state: () => ({
    selectedPatient: null as Patient | null,
  }),
  actions: {
    viewPatient(patient: Patient) {
      this.selectedPatient = patient;
    },
    clearPatient() {
      this.selectedPatient = null;
    }
  },
})