import { describe, it, expect, beforeEach } from 'vitest'
import { useNotification } from './useNotification'

describe('useNotification', () => {
  let notif: ReturnType<typeof useNotification>

  beforeEach(() => {
    notif = useNotification()
    notif.clear()
  })

  it('adds a notification via success', () => {
    notif.success('All good')
    expect(notif.notifications.value.length).toBe(1)
    expect(notif.notifications.value[0].type).toBe('success')
    expect(notif.notifications.value[0].message).toBe('All good')
  })

  it('adds a notification via error', () => {
    notif.error('Something broke')
    expect(notif.notifications.value[0].type).toBe('error')
  })

  it('adds a notification via warning', () => {
    notif.warning('Be careful')
    expect(notif.notifications.value[0].type).toBe('warning')
  })

  it('dismisses a notification', () => {
    notif.info('Test')
    const id = notif.notifications.value[0].id
    notif.dismiss(id)
    expect(notif.notifications.value.length).toBe(0)
  })
})
