import type { CheckupType } from '@/enums/checkupType';

export interface CheckupDto {
    uuid: string;
    checkupDate: string
    type: CheckupType
    medicalRecordUuid: string; 
    IllnessID?: number
}

export interface CreateCheckupDto {
    checkupDate: string;
    type: CheckupType;
    medicalRecordUuid: string;
    IllnessID?: number;
}

export interface UpdateCheckupDto {
    checkupDate: string;
    type: CheckupType;
    IllnessID?: number;
}