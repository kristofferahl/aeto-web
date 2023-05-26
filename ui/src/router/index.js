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
      path: '/resourcetemplates/:namespace?/:name?',
      name: 'resourcetemplates',
      component: () => import('../views/ResourceTemplatesView.vue')
    },
    {
      path: '/eventstreamchunks/:namespace?/:name?',
      name: 'eventstreamchunks',
      component: () => import('../views/EventStreamChunksView.vue')
    },
    {
      path: '/savingspolicies/:namespace?/:name?',
      name: 'savingspolicies',
      component: () => import('../views/SavingsPoliciesView.vue')
    },
    {
      path: '/certificates/:namespace?/:name?',
      name: 'certificates',
      component: () => import('../views/CertificatesView.vue')
    },
    {
      path: '/certificateconnectors/:namespace?/:name?',
      name: 'certificateconnectors',
      component: () => import('../views/CertificateConnectorsView.vue')
    },
    {
      path: '/hostedzones/:namespace?/:name?',
      name: 'hostedzones',
      component: () => import('../views/HostedZonesView.vue')
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue')
    }
  ]
})

export default router
