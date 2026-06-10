<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import WmsInput from '@/components/common/WmsInput.vue'
import WmsButton from '@/components/common/WmsButton.vue'
import WmsPasswordInput from '@/components/auth/WmsPasswordInput.vue'

const props = defineProps<{
  loading?: boolean
  inviteToken?: string
}>()

const emit = defineEmits<{
  submit: [data: {
    token: string
    company_name: string
    email: string
    password: string
    phone: string
    vat_number: string
  }]
}>()

const companyName = ref('')
const email = ref('')
const password = ref('')
const phone = ref('')
const vatNumber = ref('')

const touched = reactive({
  companyName: false,
  email: false,
  password: false,
  phone: false,
  vatNumber: false,
})

const companyNameError = computed(() => {
  if (!touched.companyName && !companyName.value) return ''
  if (!companyName.value) return 'Company name is required'
  return ''
})

const emailError = computed(() => {
  if (!touched.email && !email.value) return ''
  if (!email.value) return 'Email is required'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) return 'Invalid email format'
  return ''
})

const passwordError = computed(() => {
  if (!touched.password && !password.value) return ''
  if (!password.value) return 'Password is required'
  if (password.value.length < 8) return 'Password must be at least 8 characters'
  return ''
})

const passwordStrength = computed(() => {
  const val = password.value
  if (!val) return { level: 0, label: '' }
  let score = 0
  if (val.length >= 8) score++
  if (val.length >= 12) score++
  if (/[A-Z]/.test(val)) score++
  if (/[0-9]/.test(val)) score++
  if (/[^A-Za-z0-9]/.test(val)) score++
  const labels = ['', 'Weak', 'Fair', 'Good', 'Strong', 'Very Strong']
  return { level: Math.min(score, 5), label: labels[score] }
})

const isValid = computed(() => {
  return companyName.value && email.value && password.value.length >= 8 && !emailError.value
})

function onSubmit() {
  touched.companyName = true
  touched.email = true
  touched.password = true
  if (!isValid.value) return
  emit('submit', {
    token: props.inviteToken || '',
    company_name: companyName.value,
    email: email.value,
    password: password.value,
    phone: phone.value,
    vat_number: vatNumber.value,
  })
}
</script>

<template>
  <form class="wms-register-form" @submit.prevent="onSubmit">
    <h2 class="wms-register-form__title">Create Your Account</h2>
    <p class="wms-register-form__subtitle">Fill in the details to get started.</p>

    <div class="wms-register-form__fields">
      <WmsInput
        v-model="companyName"
        label="Company Name"
        placeholder="Your company name"
        :error="companyNameError"
        @blur="touched.companyName = true"
      />

      <WmsInput
        v-model="email"
        label="Email"
        type="email"
        placeholder="you@example.com"
        :error="emailError"
        @blur="touched.email = true"
      />

      <WmsPasswordInput
        v-model="password"
        label="Password"
        placeholder="Create a password"
        :error="passwordError"
        show-strength
        :strength="passwordStrength.level"
        :strength-label="passwordStrength.label"
        @blur="touched.password = true"
      />

      <WmsInput
        v-model="phone"
        label="Phone (optional)"
        type="text"
        placeholder="+1 (555) 000-0000"
      />

      <WmsInput
        v-model="vatNumber"
        label="VAT Number (optional)"
        type="text"
        placeholder="EU-VAT-XXXXXX"
      />
    </div>

    <WmsButton
      type="submit"
      variant="primary"
      size="lg"
      :loading="loading"
      :disabled="loading"
    >
      Create Account
    </WmsButton>

  </form>
</template>

<style scoped>
.wms-register-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-5);
  width: 100%;
  max-width: 28rem;
  padding: var(--spacing-8);
  background-color: var(--color-card-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
}

.wms-register-form__title {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin: 0;
}

.wms-register-form__subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.wms-register-form__fields {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

</style>
