import type { UserRole } from '@/enums/userRole';

export interface User {
  id?: number;
  uuid: string;
  firstName: string;
  lastName: string;
  oib: string;
  birthDate: string;
  gender: string;
  email: string;
  role: UserRole;
}