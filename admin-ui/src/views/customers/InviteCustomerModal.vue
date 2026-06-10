<script setup lang="ts">
import { ref, computed } from 'vue'
import WmsCard from '@/components/common/WmsCard.vue'
import WmsInput from '@/components/common/WmsInput.vue'
import WmsButton from '@/components/common/WmsButton.vue'
import { useInviteStore } from '@/stores/invites'
import { useToastStore } from '@/stores/toast'
import type { InviteUserResponse } from '@/api/types'

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const inviteStore = useInviteStore()
const toastStore = useToastStore()

const email = ref('')
const loading = ref(false)
const inviteResult = ref<InviteUserResponse | null>(null)
const copied = ref(false)

const emailError = computed(() => {
  if (!email.value) return ''
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) return 'Invalid email format'
  return ''
})

async function handleSubmit() {
  if (!email.value || emailError.value) return
  loading.value = true
  inviteResult.value = null
  copied.value = false
  try {
    const result = await inviteStore.inviteUser(email.value)
    if (typeof result === 'string') {
      // Error string from store
      toastStore.addToast(result, 'error')
      return
    }
    inviteResult.value = result
    toastStore.addToast('Invitation sent!', 'success')
  } catch (err: any) {
    toastStore.addToast(err?.message || 'Failed to send invitation', 'error')
  } finally {
    loading.value = false
  }
}

async function copyLink() {
  if (!inviteResult.value?.invite_link) return
  try {
    await navigator.clipboard.writeText(inviteResult.value.invite_link)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // Fallback for environments without clipboard API
  }
}

function handleClose() {
  email.value = ''
  inviteResult.value = null
  copied.value = false
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="invite-overlay" @click.self="handleClose">
      <WmsCard padding="lg" class="invite-card">
        <template #header>
          <h2 class="invite-card__title">Invite Customer</h2>
        </template>

        <!-- Form -->
        <form v-if="!inviteResult" @submit.prevent="handleSubmit" class="invite-card__form">
          <WmsInput
            v-model="email"
            label="Email address"
            type="email"
            placeholder="customer@company.com"
            :error="emailError"
          />

          <div class="invite-card__actions">
            <WmsButton variant="secondary" type="button" @click="handleClose">
              Cancel
            </WmsButton>
            <WmsButton variant="primary" type="submit" :loading="loading" :disabled="loading || !email">
              Send Invitation
            </WmsButton>
          </div>
        </form>

        <!-- Success -->
        <div v-else class="invite-card__result">
          <p class="invite-card__success-msg">Invitation sent successfully!</p>
          <div class="invite-card__link-row">
            <input
              :value="inviteResult.invite_link"
              class="invite-card__link-input"
              readonly
            />
            <WmsButton variant="secondary" size="sm" @click="copyLink">
              {{ copied ? 'Copied!' : 'Copy' }}
            </WmsButton>
          </div>
          <WmsButton variant="primary" @click="handleClose">
            Done
          </WmsButton>
        </div>
      </WmsCard>
    </div>
  </Teleport>
</template>

<style scoped>
.invite-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.invite-card {
  width: 100%;
  max-width: 32rem;
}

.invite-card__title {
  margin: 0;
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.invite-card__form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-5);
}

.invite-card__actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-3);
}

.invite-card__result {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.invite-card__success-msg {
  margin: 0;
  color: var(--color-success-text);
  font-weight: var(--font-medium);
}

.invite-card__link-row {
  display: flex;
  gap: var(--spacing-2);
  align-items: center;
}

.invite-card__link-input {
  flex: 1;
  padding: var(--spacing-2) var(--spacing-3);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-md);
  background-color: var(--color-input-bg);
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-family: monospace;
}
</style>
