import type { PatientDto, NewPatientDto } from '@/dtos/patientDto';
import { createNewPatientDto } from '@/dtos/patientDto';

import { formatDate } from '@/utils/formatDate';
import axios from './axios';
import type { CheckupDto, CreateCheckupDto, UpdateCheckupDto } from '@/dtos/checkupDto';
import type { CreateIllnessDto, IllnessListDto, UpdateIllnessDto } from '@/dtos/illnessDto';
import type { PrescriptionListDto, CreatePrescriptionDto, MedicationListDto } from '@/dtos/prescriptionDto';

const BASE_URL_PATIENTS = '/patients';
const BASE_URL_CHECKUPS = '/checkup';
const BASE_URL_ILLNESSES = '/illnesses';
const BASE_URL_PRESCRIPTIONS = '/prescriptions';
const BASE_URL_MEDICATIONS = '/medications';

export async function getAllPatients(): Promise<PatientDto[]> {
  const response = await axios.get<PatientDto[]>(BASE_URL_PATIENTS);
  return response.data;
}

export async function getPatientById(id: number): Promise<PatientDto> {
  const response = await axios.get<PatientDto>(`${BASE_URL_PATIENTS}/${id}`);
  return response.data;
}

export async function createPatient(newPatient: NewPatientDto): Promise<PatientDto> {
  const patientDto = createNewPatientDto(newPatient);
  patientDto.birthDate = formatDate(patientDto.birthDate);
  console.log(patientDto)
  const response = await axios.post<PatientDto>(BASE_URL_PATIENTS, patientDto);
  return response.data;
}

export async function updatePatient(id: number, patient: PatientDto): Promise<PatientDto> {
  const response = await axios.put<PatientDto>(`${BASE_URL_PATIENTS}/${id}`, patient);
  return response.data;
}

export async function deletePatient(id: number): Promise<void> {
  await axios.delete(`${BASE_URL_PATIENTS}/${id}`);
}

//CHECKUPS

export async function createCheckup(checkupDto: CreateCheckupDto): Promise<CheckupDto> {
  const response = await axios.post<CheckupDto>(BASE_URL_CHECKUPS, checkupDto);
  return response.data;
}

export async function getCheckupsForRecord(recordUuid: string): Promise<CheckupDto[]> {
    const response = await axios.get<CheckupDto[]>(`${BASE_URL_CHECKUPS}/record/${recordUuid}`);
    return response.data;
}

export async function updateCheckup(uuid: string, checkupData: UpdateCheckupDto): Promise<CheckupDto> {
  const response = await axios.put<CheckupDto>(`${BASE_URL_CHECKUPS}/${uuid}`, checkupData);
  return response.data;
}

export async function deleteCheckup(uuid: string): Promise<void> {
  await axios.delete(`${BASE_URL_CHECKUPS}/${uuid}`);
}

//ILLNESSES

export async function getIllnessesForRecord(recordUuid: string): Promise<IllnessListDto[]> {
    const response = await axios.get<IllnessListDto[]>(`${BASE_URL_ILLNESSES}/record/${recordUuid}`);
    return response.data;
}

export async function createIllness(illnessData: CreateIllnessDto): Promise<IllnessListDto> {
    const response = await axios.post<IllnessListDto>(BASE_URL_ILLNESSES, illnessData);
    return response.data;
}

export async function updateIllness(uuid: string, illnessData: UpdateIllnessDto): Promise<IllnessListDto> {
    const response = await axios.put<IllnessListDto>(`${BASE_URL_ILLNESSES}/${uuid}`, illnessData);
    return response.data;
}

export async function deleteIllness(uuid: string): Promise<void> {
    await axios.delete(`${BASE_URL_ILLNESSES}/${uuid}`);
}

//PRESCRIPTIONS

export async function getPrescriptionsForIllness(illnessId: number): Promise<PrescriptionListDto[]> {
    const response = await axios.get<PrescriptionListDto[]>(`${BASE_URL_PRESCRIPTIONS}/illness/${illnessId}`);
    return response.data;
}

export async function createPrescription(prescriptionData: CreatePrescriptionDto): Promise<PrescriptionListDto> {
    const response = await axios.post<PrescriptionListDto>(BASE_URL_PRESCRIPTIONS, prescriptionData);
    return response.data;
}

export async function deletePrescription(uuid: string): Promise<void> {
    await axios.delete(`${BASE_URL_PRESCRIPTIONS}/${uuid}`);
}

//MEDICATIONS

export async function getAllMedications(): Promise<MedicationListDto[]> {
  const response = await axios.get<MedicationListDto[]>(BASE_URL_MEDICATIONS);
  return response.data;
}