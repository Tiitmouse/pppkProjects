import { UserRole } from "@/enums/userRole";

export interface NavigationGroup {
    Name: string;
    Links: NavigationLink[]
  }
  
export interface NavigationLink {
    Name: string;
    AllowRoles: UserRole[];
    Route: string;
  }
  
 export const navigationLinks: NavigationGroup[] = ([
    {
      Name: "Patients",
      Links: [
        {
          Name: "Overview",
          AllowRoles: [UserRole.Doctor],
          Route: '/patient-list'
        }
      ]
    },
    {
      Name: "Actions",
      Links: [
         {
          Name: "Creation",
          AllowRoles: [UserRole.SuperAdmin, UserRole.Doctor],
          Route: '/new-user'
        },
      ]
    },
  ]);