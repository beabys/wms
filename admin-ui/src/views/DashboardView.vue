<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
</script>

<template>
  <div class="dashboard">
    <header class="dashboard__header">
      <h1>Dashboard</h1>
    </header>
    <main class="dashboard__content">
      <p class="dashboard__welcome">
        Welcome{{ auth.user?.name ? `, ${auth.user.name}` : '' }}!
      </p>
      <p class="dashboard__subtitle">
        Use the sidebar to navigate between sections.
      </p>
    </main>
  </div>
</template>

<style scoped>
.dashboard {
  padding: var(--spacing-6);
  max-width: 80rem;
  margin: 0 auto;
}

.dashboard__header {
  margin-bottom: var(--spacing-6);
}

.dashboard__header h1 {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin: 0;
}

.dashboard__content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.dashboard__welcome {
  font-size: var(--text-lg);
  color: var(--color-text-primary);
  margin: 0;
}

.dashboard__subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin: 0;
}
</style>
