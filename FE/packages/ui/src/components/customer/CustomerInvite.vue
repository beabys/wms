<script setup lang="ts">
import { ref } from 'vue'
defineOptions({ name: 'CustomerInvite' })

const props = defineProps<{
  inviteLink: string
  expiresAt: string
}>()

const copied = ref(false)

function copyToClipboard(): void {
  navigator.clipboard.writeText(props.inviteLink).then(() => {
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  })
}
</script>

<template>
  <div class="bg-white border rounded-lg p-6 space-y-4">
    <h3 class="text-lg font-semibold text-gray-800">Invite Customer</h3>
    <p class="text-sm text-gray-500">
      Share this link with your customer to let them sign up. Expires {{ expiresAt }}.
    </p>
    <div class="flex items-center gap-2">
      <input
        :value="inviteLink"
        readonly
        class="flex-1 px-3 py-2 border rounded text-sm bg-gray-50 text-gray-600"
      />
      <button
        class="px-4 py-2 bg-blue-600 text-white rounded text-sm font-medium hover:bg-blue-700 flex items-center gap-1"
        @click="copyToClipboard"
      >
        <span v-if="copied">&#10003;</span>
        {{ copied ? 'Copied' : 'Copy' }}
      </button>
    </div>
  </div>
</template>
