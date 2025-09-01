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
  // Routes previously generated from "Akcije" group
  {
    path: '/new-user',
    name: 'new-user',
    component: () => import(/* @vite-ignore */ '@/pages/new-user.vue'),
    meta: {
      allowedRoles: [UserRole.SuperAdmin],
      layout: 'default',
      breadcrumbName: "Novi korisnik"
    },
  },
  {
    path: '/users',
    name: 'users',
    component: () => import(/* @vite-ignore */ '@/pages/users.vue'),
    meta: {
      allowedRoles: [UserRole.SuperAdmin],
      layout: 'default',
      breadcrumbName: "Korisnici"
    },
  },
  {
    path: '/details',
    name: 'details',
    component: () => import('@/pages/details.vue'),
    meta: {
      allowedRoles: [],
      layout: 'default',
      breadcrumbName: "Tvoji podaci"
    },
  }

]
