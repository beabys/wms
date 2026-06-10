<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'
import WmsLayout from '@/components/layout/WmsLayout.vue'
import WmsToast from '@/components/layout/WmsToast.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
useTheme().init()
</script>

<template>
  <WmsToast />
  <WmsLayout v-if="auth.isAuthenticated">
    <router-view />
  </WmsLayout>
  <router-view v-else />
</template>
