import type { CheckupType } from '@/enums/checkupType';

export interface CheckupDto {
    uuid: string;
    checkupDate: string;
    type: CheckupType;
    medicalRecordUuid: string;
    illnessId?: number;
}

export interface CreateCheckupDto {
    checkupDate: string;
    type: CheckupType;
    medicalRecordUuid: string;
    illnessId?: number;
}

export interface UpdateCheckupDto {
    checkupDate: string;
    type: CheckupType;
    illnessId?: number;
}