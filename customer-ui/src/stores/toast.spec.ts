import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useToastStore } from './toast'

describe('useToastStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('starts with no toasts', () => {
    const store = useToastStore()
    expect(store.toasts).toEqual([])
  })

  it('adds a toast with default type info', () => {
    const store = useToastStore()
    store.addToast('Hello')
    expect(store.toasts).toHaveLength(1)
    expect(store.toasts[0].message).toBe('Hello')
    expect(store.toasts[0].type).toBe('info')
    expect(store.toasts[0].id).toBe(0)
  })

  it('adds a toast with specified type', () => {
    const store = useToastStore()
    store.addToast('Error!', 'error')
    expect(store.toasts[0].type).toBe('error')
  })

  it('removes a toast by id', () => {
    const store = useToastStore()
    store.addToast('First')
    store.addToast('Second')
    store.removeToast(0)
    expect(store.toasts).toHaveLength(1)
    expect(store.toasts[0].message).toBe('Second')
  })

  it('auto-dismisses a toast after duration', () => {
    const store = useToastStore()
    store.addToast('Auto dismiss', 'info', 4000)
    expect(store.toasts).toHaveLength(1)

    vi.advanceTimersByTime(4000)
    expect(store.toasts).toHaveLength(0)
  })

  it('does not auto-dismiss when duration is 0', () => {
    const store = useToastStore()
    store.addToast('Sticky', 'info', 0)
    vi.advanceTimersByTime(99999)
    expect(store.toasts).toHaveLength(1)
  })

  it('clears all toasts', () => {
    const store = useToastStore()
    store.addToast('A')
    store.addToast('B')
    store.addToast('C')
    expect(store.toasts).toHaveLength(3)
    store.clearAll()
    expect(store.toasts).toHaveLength(0)
  })

  it('increments ids for each toast', () => {
    const store = useToastStore()
    store.addToast('A')
    store.addToast('B')
    store.addToast('C')
    expect(store.toasts[0].id).toBe(0)
    expect(store.toasts[1].id).toBe(1)
    expect(store.toasts[2].id).toBe(2)
  })
})
