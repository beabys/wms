import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/customer-approvals',
      name: 'CustomerApprovals',
      component: () => import('../views/CustomerApprovalsView.vue')
    },
    {
      path: '/inbound-queue',
      name: 'InboundQueue',
      component: () => import('../views/InboundQueueView.vue')
    },
    {
      path: '/inbound/:id',
      name: 'InboundDetail',
      component: () => import('../views/InboundQueueView.vue')
    },
    {
      path: '/inventory',
      name: 'Inventory',
      component: () => import('../views/InventoryView.vue')
    },
    {
      path: '/inventory/products',
      name: 'Products',
      component: () => import('../views/ProductsView.vue')
    },
    {
      path: '/inventory/bin-locations',
      name: 'BinLocations',
      component: () => import('../views/BinLocationsView.vue')
    },
    {
      path: '/stock',
      name: 'Stock',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/orders',
      name: 'Orders',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/pick-lists',
      name: 'PickLists',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/returns',
      name: 'Returns',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/billing',
      name: 'Billing',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/reports',
      name: 'Reports',
      component: () => import('../views/DashboardView.vue')
    }
  ]
})

export default router
