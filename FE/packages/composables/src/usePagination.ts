import { ref, computed } from 'vue'

export function usePagination(initialPage = 1, initialLimit = 10) {
  const page = ref(initialPage)
  const limit = ref(initialLimit)
  const total = ref(0)

  const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit.value)))
  const hasNext = computed(() => page.value < totalPages.value)
  const hasPrev = computed(() => page.value > 1)

  function setPage(p: number): void {
    page.value = Math.max(1, Math.min(p, totalPages.value))
  }

  function nextPage(): void {
    if (hasNext.value) page.value++
  }

  function prevPage(): void {
    if (hasPrev.value) page.value--
  }

  function setTotal(t: number): void {
    total.value = t
    if (page.value > totalPages.value) {
      page.value = totalPages.value
    }
  }

  return { page, limit, total, totalPages, hasNext, hasPrev, setPage, nextPage, prevPage, setTotal }
}
