<script setup lang="ts">
import WmsSidebar from './WmsSidebar.vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { logout } = useAuth()

async function handleLogout() {
  await logout()
  router.push('/login')
}
</script>

<template>
  <div class="wms-layout">
    <WmsSidebar @logout="handleLogout" />
    <main class="wms-layout__content">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.wms-layout {
  display: flex;
  min-height: 100vh;
  background-color: var(--color-bg-primary);
}

.wms-layout__content {
  flex: 1;
  margin-left: 240px;
  overflow-y: auto;
}
</style>
