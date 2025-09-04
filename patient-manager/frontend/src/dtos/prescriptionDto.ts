export interface MedicationListDto {
    uuid: string;
    name: string;
}

export interface PrescriptionListDto {
    uuid: string;
    issuedAt: string;
    medications: MedicationListDto[];
}

export interface CreatePrescriptionDto {
    issuedAt: string;
    illnessId: number;
    medicationUuids: string[];
}