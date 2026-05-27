<script setup lang="ts">
import { computed } from 'vue'
defineOptions({ name: 'WmsPagination' })

const props = defineProps<{
  page: number
  total: number
  limit: number
}>()

const emit = defineEmits<{
  change: [page: number]
}>()

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)))

const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, props.page - 2)
  const end = Math.min(totalPages.value, props.page + 2)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

function goTo(page: number): void {
  if (page >= 1 && page <= totalPages.value && page !== props.page) {
    emit('change', page)
  }
}
</script>

<template>
  <nav class="inline-flex items-center gap-1">
    <button
      class="px-2 py-1 text-sm rounded border hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      :disabled="page <= 1"
      @click="goTo(1)"
    >
      &laquo;
    </button>
    <button
      class="px-2 py-1 text-sm rounded border hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      :disabled="page <= 1"
      @click="goTo(page - 1)"
    >
      &lsaquo;
    </button>
    <button
      v-for="p in visiblePages"
      :key="p"
      class="px-3 py-1 text-sm rounded border"
      :class="p === page ? 'bg-blue-600 text-white border-blue-600' : 'hover:bg-gray-50'"
      @click="goTo(p)"
    >
      {{ p }}
    </button>
    <button
      class="px-2 py-1 text-sm rounded border hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      :disabled="page >= totalPages"
      @click="goTo(page + 1)"
    >
      &rsaquo;
    </button>
    <button
      class="px-2 py-1 text-sm rounded border hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      :disabled="page >= totalPages"
      @click="goTo(totalPages)"
    >
      &raquo;
    </button>
    <span class="ml-2 text-sm text-gray-500">
      Page {{ page }} of {{ totalPages }}
    </span>
  </nav>
</template>
