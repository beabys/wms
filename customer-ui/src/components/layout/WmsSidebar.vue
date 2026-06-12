<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useCustomerStore } from '@/stores/customer'
import WmsThemeToggle from './WmsThemeToggle.vue'

const emit = defineEmits<{
  logout: []
}>()

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const customerStore = useCustomerStore()

const navLinks = [
  { label: 'Dashboard', path: '/dashboard' },
]

function isActive(path: string): boolean {
  return route.path === path
}

function navigate(path: string) {
  router.push(path)
}
</script>

<template>
  <aside class="wms-sidebar">
    <div class="wms-sidebar__top">
      <div class="wms-sidebar__brand">
        <span class="wms-sidebar__logo">WMS</span>
        <span class="wms-sidebar__title">Customer</span>
      </div>
    </div>

    <nav v-if="customerStore.isApproved" class="wms-sidebar__nav">
      <button
        v-for="link in navLinks"
        :key="link.path"
        :class="['wms-sidebar__link', { 'wms-sidebar__link--active': isActive(link.path) }]"
        @click="navigate(link.path)"
      >
        {{ link.label }}
      </button>
    </nav>
    <div v-else class="wms-sidebar__spacer" />

    <div class="wms-sidebar__bottom">
      <div v-if="auth.user" class="wms-sidebar__user" :title="auth.user.email">
        <span class="wms-sidebar__user-email">{{ auth.user.email }}</span>
      </div>
      <div class="wms-sidebar__actions">
        <WmsThemeToggle />
        <button class="wms-sidebar__logout wms-sidebar__logout--ghost" @click="emit('logout')">
          Logout
        </button>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.wms-sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 240px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--color-card-bg);
  border-right: 1px solid var(--color-border);
  z-index: 100;
}

.wms-sidebar__top {
  padding: var(--spacing-4) var(--spacing-5);
  border-bottom: 1px solid var(--color-border);
}

.wms-sidebar__brand {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.wms-sidebar__logo {
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-primary);
}

.wms-sidebar__title {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.wms-sidebar__nav {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-2) 0;
}

.wms-sidebar__spacer {
  flex: 1;
}

.wms-sidebar__link {
  display: block;
  width: 100%;
  padding: var(--spacing-2-5) var(--spacing-5);
  border: none;
  background: transparent;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  text-align: left;
  cursor: pointer;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.wms-sidebar__link:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.wms-sidebar__link--active {
  background-color: var(--color-primary-bg);
  color: var(--color-primary);
  font-weight: var(--font-semibold);
}

.wms-sidebar__bottom {
  padding: var(--spacing-4) var(--spacing-5);
  border-top: 1px solid var(--color-border);
}

.wms-sidebar__user {
  margin-bottom: var(--spacing-3);
  overflow: hidden;
}

.wms-sidebar__user-email {
  display: block;
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.wms-sidebar__actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.wms-sidebar__logout {
  flex: 1;
  padding: var(--spacing-1-5) var(--spacing-3);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  cursor: pointer;
  transition: background-color 0.15s ease, color 0.15s ease;
  white-space: nowrap;
}

.wms-sidebar__logout--ghost {
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
}

.wms-sidebar__logout--ghost:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}
</style>
