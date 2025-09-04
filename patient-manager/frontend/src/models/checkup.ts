import type { CheckupType } from "@/enums/checkupType"


export interface Checkup {
    checkupDate: string
    type: CheckupType
    medicalRecordID: number
    IllnessID: number
}