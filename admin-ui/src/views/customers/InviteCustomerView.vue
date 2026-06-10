<script setup lang="ts">
import { useRouter } from 'vue-router'
import InviteCustomerModal from './InviteCustomerModal.vue'

const router = useRouter()

function handleClose() {
  router.push('/customers')
}
</script>

<template>
  <InviteCustomerModal :visible="true" @close="handleClose" />
</template>
