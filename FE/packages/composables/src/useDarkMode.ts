import { ref, watch } from 'vue'

const STORAGE_KEY = 'wms_dark_mode'

function getInitialValue(): boolean {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored !== null) return stored === 'true'
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

const isDark = ref(getInitialValue())

function apply(value: boolean) {
  document.documentElement.classList.toggle('dark', value)
}

// Apply on load
apply(isDark.value)

// Persist changes (sync so tests can verify immediately)
watch(isDark, (val) => {
  localStorage.setItem(STORAGE_KEY, String(val))
  apply(val)
}, { flush: 'sync' })

export function useDarkMode() {
  function toggle() { isDark.value = !isDark.value }
  return { isDark, toggle }
}
