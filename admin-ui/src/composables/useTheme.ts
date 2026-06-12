import { ref } from 'vue'

const THEME_KEY = 'wms_admin_theme'
const isDark = ref(false)

export function useTheme() {
  function init() {
    const saved = localStorage.getItem(THEME_KEY)
    if (saved === 'dark') {
      isDark.value = true
      document.documentElement.setAttribute('data-theme', 'dark')
    } else {
      isDark.value = false
      document.documentElement.setAttribute('data-theme', 'light')
    }
  }

  function toggle() {
    isDark.value = !isDark.value
    const theme = isDark.value ? 'dark' : 'light'
    document.documentElement.setAttribute('data-theme', theme)
    localStorage.setItem(THEME_KEY, theme)
  }

  function setTheme(theme: 'light' | 'dark') {
    isDark.value = theme === 'dark'
    document.documentElement.setAttribute('data-theme', theme)
    localStorage.setItem(THEME_KEY, theme)
  }

  return {
    isDark,
    init,
    toggle,
    setTheme,
  }
}
