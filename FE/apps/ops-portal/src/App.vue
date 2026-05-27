<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useDarkMode } from '@wms/composables'

const router = useRouter()
const auth = useAuthStore()
const { isDark, toggle: toggleDark } = useDarkMode()


const navLinks = [
  { label: 'Dashboard', path: '/dashboard' },
  { label: 'Customer Approvals', path: '/customer-approvals' },
  { label: 'Inbound Queue', path: '/inbound-queue' },
  { label: 'Inventory', path: '/inventory' },
  { label: 'Products', path: '/inventory/products' },
  { label: 'Bin Locations', path: '/inventory/bin-locations' },
  { label: 'Stock', path: '/stock' },
  { label: 'Orders', path: '/orders' },
  { label: 'Pick Lists', path: '/pick-lists' },
  { label: 'Returns', path: '/returns' },
  { label: 'Billing', path: '/billing' },
  { label: 'Reports', path: '/reports' }
]

function handleLogout(): void {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div v-if="auth.isAuthenticated" class="flex h-screen bg-gray-100 dark:bg-gray-950">
    <!-- Sidebar -->
    <aside class="w-60 bg-gray-900 dark:bg-gray-800 text-white flex flex-col shrink-0">
      <div class="px-4 py-5 border-b border-gray-700 dark:border-gray-600">
        <h1 class="text-lg font-bold">WMS Operations</h1>
      </div>
      <nav class="flex-1 p-3 space-y-1">
        <router-link
          v-for="link in navLinks"
          :key="link.path"
          :to="link.path"
          class="block px-3 py-2 rounded text-sm hover:bg-gray-700 dark:hover:bg-gray-600"
          active-class="bg-blue-600"
        >
          {{ link.label }}
        </router-link>
      </nav>
    </aside>

    <!-- Main -->
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white dark:bg-gray-900 border-b dark:border-gray-700 flex items-center justify-between px-6 shrink-0">
        <span class="text-sm text-gray-500 dark:text-gray-400">Operations Portal</span>
        <div class="flex items-center gap-4">
          <button
            class="text-sm px-2 py-1 rounded border dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="toggleDark"
            :title="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
          >
            <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"/></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/></svg>
          </button>
          <span class="text-sm text-gray-700 dark:text-gray-300">{{ auth.user?.name }}</span>
          <button
            class="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
            @click="handleLogout"
          >
            Logout
          </button>
        </div>
      </header>
      <main class="flex-1 overflow-auto bg-gray-100 dark:bg-gray-950">
        <router-view />
      </main>
    </div>
  </div>

  <router-view v-else />
</template>
