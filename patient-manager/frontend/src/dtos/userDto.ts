export interface UserDto {
    uuid: string;
    firstName: string;
    lastName: string;
    oib: string;
    residence: string;
    birthDate: string;
    email: string;
    role: string;
}

export interface DoctorDto {
    firstName: string;
    lastName: string;
}