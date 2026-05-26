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
      path: '/sign-up',
      name: 'SignUp',
      component: () => import('../views/SignUpView.vue')
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/inbound/new',
      name: 'InboundForm',
      component: () => import('../views/InboundFormView.vue')
    },
    {
      path: '/inbound/list',
      name: 'InboundList',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/my-stock',
      name: 'MyStock',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/orders',
      name: 'Orders',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/returns',
      name: 'Returns',
      component: () => import('../views/DashboardView.vue')
    },
    {
      path: '/ledger',
      name: 'Ledger',
      component: () => import('../views/DashboardView.vue')
    }
  ]
})

export default router
