import type { CheckupType } from '@/enums/checkupType';

export interface CheckupDto {
    checkupDate: string
    type: CheckupType
    medicalRecordUuid: string; 
    IllnessID?: number
}