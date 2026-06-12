<script setup lang="ts">
import { onMounted } from 'vue'
import { useCustomerStore } from '@/stores/customer'

const customerStore = useCustomerStore()

onMounted(() => {
  customerStore.fetchMyCustomer()
})
</script>

<template>
  <div class="dashboard">
    <header class="dashboard__header">
      <h1>Dashboard</h1>
    </header>
    <main class="dashboard__content">
      <!-- Pending approval -->
      <div v-if="customerStore.approvalStatus === 'pending'" class="dashboard__status dashboard__status--pending">
        <h2>Account Pending Approval</h2>
        <p>Your account is awaiting approval from an administrator. You will be able to access all features once your account is approved.</p>
      </div>

      <!-- Rejected -->
      <div v-else-if="customerStore.approvalStatus === 'rejected'" class="dashboard__status dashboard__status--rejected">
        <h2>Account Rejected</h2>
        <p>Your account has been rejected. Please contact support for more information.</p>
      </div>

      <!-- Suspended -->
      <div v-else-if="customerStore.approvalStatus === 'suspended'" class="dashboard__status dashboard__status--suspended">
        <h2>Account Suspended</h2>
        <p>Your account has been suspended. Please contact support.</p>
      </div>

      <!-- Active / Approved -->
      <div v-else-if="customerStore.isApproved">
        <p>Welcome to the WMS Customer Portal.</p>
      </div>

      <!-- Loading -->
      <div v-else>
        <p>Loading...</p>
      </div>
    </main>
  </div>
</template>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: var(--color-bg-primary);
}

.dashboard__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-4) var(--spacing-6);
  border-bottom: 1px solid var(--color-border);
  background-color: var(--color-card-bg);
}

.dashboard__header h1 {
  font-size: var(--text-xl);
}

.dashboard__content {
  padding: var(--spacing-6);
}

.dashboard__status {
  padding: var(--spacing-8);
  border-radius: var(--radius-lg);
  text-align: center;
}

.dashboard__status--pending {
  background-color: var(--color-warning-bg);
  color: var(--color-warning-text);
  border: 1px solid var(--color-warning);
}

.dashboard__status--rejected {
  background-color: var(--color-error-bg);
  color: var(--color-error-text);
  border: 1px solid var(--color-error);
}

.dashboard__status--suspended {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.dashboard__status h2 {
  margin: 0 0 var(--spacing-2);
  font-size: var(--text-lg);
}

.dashboard__status p {
  margin: 0;
  font-size: var(--text-sm);
}
</style>
