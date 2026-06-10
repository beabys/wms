<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import WmsRegisterForm from '@/components/auth/WmsRegisterForm.vue'
import WmsThemeToggle from '@/components/layout/WmsThemeToggle.vue'
import { useAuth } from '@/composables/useAuth'
import { useToastStore } from '@/stores/toast'
import type { RegisterRequest } from '@/api/types'

const router = useRouter()
const route = useRoute()
const { register } = useAuth()
const toastStore = useToastStore()
const inviteToken = ref('')
const loading = ref(false)
const noToken = computed(() => !inviteToken.value)

onMounted(() => {
  const token = route.query.token as string
  if (token) {
    inviteToken.value = token
  }
})

async function handleSubmit(data: RegisterRequest) {
  loading.value = true
  try {
    await register(data)
    toastStore.addToast('Account created! You can now sign in.', 'success')
    router.push('/login?registered=true')
  } catch (err: any) {
    const msg = err?.message || 'Registration failed'
    if (msg.includes('missing required fields: token')) {
      toastStore.addToast('No invite token found. Please use the invite link from your admin.', 'warning', 6000)
    } else {
      toastStore.addToast(msg, 'error')
    }
  } finally {
    loading.value = false
  }
}
</script> 

<template>
  <div class="register-view">
    <div class="register-view__toggle">
      <WmsThemeToggle />
    </div>
    <div class="register-view__card">
      <div v-if="noToken" class="register-view__no-token">
        Registration requires an invite link. Ask your admin to send you one.
      </div>
      <WmsRegisterForm
        :invite-token="inviteToken"
        :loading="loading"
        @submit="handleSubmit"
      />
      <p class="register-view__login">
        Already have an account?
        <router-link to="/login">Sign in</router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
.register-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: var(--spacing-4);
  background-color: var(--color-bg-secondary);
}

.register-view__toggle {
  position: fixed;
  top: var(--spacing-4);
  right: var(--spacing-4);
}

.register-view__card {
  width: 100%;
  max-width: 32rem;
}

.register-view__login {
  margin-top: var(--spacing-4);
  text-align: center;
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.register-view__login a {
  color: var(--color-primary);
  font-weight: var(--font-medium);
}

.register-view__login a:hover {
  text-decoration: underline;
}

.register-view__no-token {
  padding: var(--spacing-3) var(--spacing-4);
  margin-bottom: var(--spacing-4);
  background-color: var(--color-warning-bg);
  color: var(--color-warning-text);
  border: 1px solid var(--color-warning);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  text-align: center;
}
</style>
