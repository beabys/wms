import { ref } from 'vue'

export type NotificationType = 'success' | 'error' | 'warning' | 'info'

export interface Notification {
  id: number
  message: string
  type: NotificationType
}

const notifications = ref<Notification[]>([])
let nextId = 0

export function useNotification() {
  function push(message: string, type: NotificationType): void {
    const id = nextId++
    notifications.value.push({ id, message, type })
    setTimeout(() => {
      notifications.value = notifications.value.filter(n => n.id !== id)
    }, 4000)
  }

  function success(message: string): void { push(message, 'success') }
  function error(message: string): void { push(message, 'error') }
  function warning(message: string): void { push(message, 'warning') }
  function info(message: string): void { push(message, 'info') }

  function dismiss(id: number): void {
    notifications.value = notifications.value.filter(n => n.id !== id)
  }

  function clear(): void {
    notifications.value = []
    nextId = 0
  }

  return { notifications, push, success, error, warning, info, dismiss, clear }
}
