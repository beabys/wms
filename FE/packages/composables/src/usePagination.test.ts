import { describe, it, expect } from 'vitest'
import { usePagination } from './usePagination'

describe('usePagination', () => {
  it('starts at page 1 with defaults', () => {
    const pag = usePagination()
    expect(pag.page.value).toBe(1)
    expect(pag.limit.value).toBe(10)
    expect(pag.total.value).toBe(0)
  })

  it('computes totalPages correctly', () => {
    const pag = usePagination(1, 10)
    pag.setTotal(25)
    expect(pag.totalPages.value).toBe(3)
  })

  it('setPage clamps to valid range', () => {
    const pag = usePagination(1, 10)
    pag.setTotal(25)
    pag.setPage(5)
    expect(pag.page.value).toBe(3) // clamped to totalPages
    pag.setPage(1)
    expect(pag.page.value).toBe(1)
  })

  it('nextPage and prevPage work', () => {
    const pag = usePagination(1, 10)
    pag.setTotal(30)
    pag.nextPage()
    expect(pag.page.value).toBe(2)
    pag.prevPage()
    expect(pag.page.value).toBe(1)
  })

  it('hasNext and hasPrev flags', () => {
    const pag = usePagination(1, 10)
    pag.setTotal(30)
    expect(pag.hasNext.value).toBe(true)
    expect(pag.hasPrev.value).toBe(false)
    pag.setPage(3)
    expect(pag.hasNext.value).toBe(false)
    expect(pag.hasPrev.value).toBe(true)
  })
})
