import { UserRole } from "@/enums/userRole";

export const routes = [
  {
    path: "/",
    name: "index",
    component: () => import('@/pages/index.vue'),
    meta: {
      layout: 'default',
      allowedRoles: []
    },
  },
  {
    path: "/login",
    name: "login",
    component: () => import('@/pages/login.vue'),
    meta: {
      layout: 'loginLayout',
    },
  },  
  {
    path: "/patient-list",
    name: "patient-list",
    component: () => import('@/pages/patient-list.vue'),
    meta: {
      layout: 'default',
    },
  },  
  {
    path: "/patient-details",
    name: "patient-details",
    component: () => import('@/pages/patient-details.vue'),
    meta: {
      layout: 'default',
    },
  },  
  {
    path: "/checkup",
    name: "checkup",
    component: () => import('@/pages/checkup.vue'),
    meta: {
      layout: 'default',
    },
  },
    {
    path: "/perscription",
    name: "perscription",
    component: () => import('@/pages/perscription.vue'),
    meta: {
      layout: 'default',
    },
  },
  {
    path: '/new-user',
    name: 'new-user',
    component: () => import(/* @vite-ignore */ '@/pages/new-user.vue'),
    meta: {
      allowedRoles: [UserRole.SuperAdmin, UserRole.Doctor],
      layout: 'default',
      breadcrumbName: "New user"
    },
  },
  {
    path: '/details',
    name: 'details',
    component: () => import('@/pages/details.vue'),
    meta: {
      allowedRoles: [],
      layout: 'default',
      breadcrumbName: "My data"
    },
  }

]
