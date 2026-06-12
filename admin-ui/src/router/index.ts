import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
  },
  {
    path: '/users',
    name: 'Users',
    component: () => import('@/views/users/UserListView.vue'),
  },
  {
    path: '/users/create',
    name: 'CreateUser',
    component: () => import('@/views/users/CreateUserView.vue'),
  },
  {
    path: '/customers',
    name: 'Customers',
    component: () => import('@/views/customers/CustomerListView.vue'),
  },
  {
    path: '/customers/:id',
    name: 'CustomerDetail',
    component: () => import('@/views/customers/CustomerDetailView.vue'),
  },
  {
    path: '/customers/invite',
    name: 'InviteList',
    component: () => import('@/views/customers/InviteListView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.name !== 'Login' && !auth.isAuthenticated) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
