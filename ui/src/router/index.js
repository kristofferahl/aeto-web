import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: { path: '/dashboard' }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView
    },
    {
      path: '/tenants/:namespace?/:name?',
      name: 'tenants',
      component: () => import('../views/TenantsView.vue')
    },
    {
      path: '/blueprints/:namespace?/:name?',
      name: 'blueprints',
      component: () => import('../views/BlueprintsView.vue')
    },
    {
      path: '/resourcesets/:namespace?/:name?',
      name: 'resourcesets',
      component: () => import('../views/ResourceSetsView.vue')
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue')
    }
  ]
})

export default router
