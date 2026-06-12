<script setup lang="ts">
import { useToastStore } from '@/stores/toast'

const store = useToastStore()

function getToastClass(type: string): string {
  return `wms-toast__item wms-toast__item--${type}`
}
</script>

<template>
  <div class="wms-toast">
    <TransitionGroup name="wms-toast-slide">
      <div
        v-for="toast in store.toasts"
        :key="toast.id"
        :class="getToastClass(toast.type)"
      >
        <span class="wms-toast__message">{{ toast.message }}</span>
        <button
          class="wms-toast__close"
          @click="store.removeToast(toast.id)"
          aria-label="Close notification"
        >
          &times;
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.wms-toast {
  position: fixed;
  top: var(--spacing-4);
  right: var(--spacing-4);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
  max-width: 24rem;
  pointer-events: none;
}

.wms-toast__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  padding: var(--spacing-3) var(--spacing-4);
  border-radius: var(--radius-md);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  pointer-events: auto;
  font-size: var(--text-sm);
  line-height: var(--leading-tight);
}

.wms-toast__message {
  flex: 1;
}

.wms-toast__close {
  flex-shrink: 0;
  width: 1.25rem;
  height: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
  font-size: 1.125rem;
  line-height: 1;
  opacity: 0.7;
  border-radius: var(--radius-sm);
}

.wms-toast__close:hover {
  opacity: 1;
}

/* Type variants */
.wms-toast__item--success {
  background-color: var(--color-success-bg);
  color: var(--color-success-text);
  border-left: 4px solid var(--color-success);
}

.wms-toast__item--error {
  background-color: var(--color-error-bg);
  color: var(--color-error-text);
  border-left: 4px solid var(--color-error);
}

.wms-toast__item--warning {
  background-color: var(--color-warning-bg);
  color: var(--color-warning-text);
  border-left: 4px solid var(--color-warning);
}

.wms-toast__item--info {
  background-color: var(--color-info-bg);
  color: var(--color-info-text);
  border-left: 4px solid var(--color-info);
}

/* Slide transition */
.wms-toast-slide-enter-active {
  transition: all 0.3s ease-out;
}

.wms-toast-slide-leave-active {
  transition: all 0.2s ease-in;
}

.wms-toast-slide-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.wms-toast-slide-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
