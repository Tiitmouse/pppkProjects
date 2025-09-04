import { defineStore } from 'pinia'

export interface Patient {
  id: number;
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string;
  gender: string;
  medicalRecordUuid: string
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