export interface IllnessListDto {
    uuid: string;
    name: string;
    startDate: string;
    endDate?: string;
}

export interface CreateIllnessDto {
    name: string;
    startDate: string;
    endDate?: string;
    medicalRecordUuid: string;
}

export interface UpdateIllnessDto {
    name: string;
    startDate: string;
    endDate?: string;
}