import type { PatientDto, NewPatientDto } from '@/dtos/patientDto';
import { createNewPatientDto } from '@/dtos/patientDto';

import { formatDate } from '@/utils/formatDate';
import axios from './axios';

const BASE_URL = '/patients';

export async function getAllPatients(): Promise<PatientDto[]> {
  const response = await axios.get<PatientDto[]>(BASE_URL);
  return response.data;
}

export async function getPatientById(id: number): Promise<PatientDto> {
  const response = await axios.get<PatientDto>(`${BASE_URL}/${id}`);
  return response.data;
}

export async function createPatient(newPatient: NewPatientDto): Promise<PatientDto> {
  const patientDto = createNewPatientDto(newPatient);
  patientDto.birthDate = formatDate(patientDto.birthDate);
  console.log(patientDto)
  const response = await axios.post<PatientDto>(BASE_URL, patientDto);
  return response.data;
}

export async function updatePatient(id: number, patient: PatientDto): Promise<PatientDto> {
  const response = await axios.put<PatientDto>(`${BASE_URL}/${id}`, patient);
  return response.data;
}

export async function deletePatient(id: number): Promise<void> {
  await axios.delete(`${BASE_URL}/${id}`);
}