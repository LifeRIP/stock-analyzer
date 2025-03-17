// frontend/src/components/__tests__/StockTable.spec.ts
import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import StockTable from '../StockTable.vue'
import { createPinia, setActivePinia } from 'pinia'

describe('StockTable', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const mockStocks = [
    {
      id: '1',
      ticker: 'AAPL',
      company: 'Apple Inc.',
      brokerage: 'Morgan Stanley',
      action: 'Buy',
      rating_from: 'Hold',
      rating_to: 'Buy',
      target_from: '$150',
      target_to: '$180',
      time: '2023-07-20T10:00:00Z',
      created_at: '2023-07-19T10:00:00Z',
      updated_at: '2023-07-20T10:00:00Z'
    }
  ]

  // Test initial rendering
  it('renders properly with props', () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })
    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.text()).toContain('AAPL')
  })

  // Test loading state
  it('shows loading spinner when isLoading is true', () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: [],
        isLoading: true,
        error: null
      }
    })
    expect(wrapper.find('.animate-spin').exists()).toBe(true)
  })

  // Test error state
  it('shows error message when error prop is provided', () => {
    const errorMessage = 'Failed to load stocks'
    const wrapper = mount(StockTable, {
      props: {
        stocks: [],
        isLoading: false,
        error: errorMessage
      }
    })
    expect(wrapper.text()).toContain(errorMessage)
  })

  // Test empty state
  it('shows empty message when no stocks', () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: [],
        isLoading: false,
        error: null
      }
    })
    expect(wrapper.text()).toContain('No stocks found')
  })

  // Test search functionality
  it('emits search event when input changes', async () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })
    const input = wrapper.find('input[type="text"]')
    await input.setValue('AAPL')
    expect(wrapper.emitted('search')?.[0]).toEqual(['AAPL'])
  })

  // Test page size changes
  it('emits page-size-change event when select changes', async () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })
    const select = wrapper.find('select')
    await select.setValue(50)
    expect(wrapper.emitted('page-size-change')?.[0]).toEqual([50])
  })

  // Test stock click
  it('emits stock-click event when row is clicked', async () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })

    await wrapper.find('td:nth-child(1)').trigger('click')
    expect(wrapper.emitted('stock-click')?.[0]).toEqual(['AAPL'])
  })

  // Test date formatting
  it('formats date correctly', () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })
    const formattedDate = wrapper.find('td:last-child').text()
    expect(formattedDate).toBeTruthy()
  })

  // Test rating class
  it('applies correct rating class', () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })
    // td:nth-child(5) > div
    const ratingCell = wrapper.find('td:nth-child(5) > div')
    expect(ratingCell.classes()).toContain('text-green-600')
  })

  // Test target class
  it('applies correct target class', () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })
    const targetCell = wrapper.find('td:nth-child(6) > div')
    expect(targetCell.classes()).toContain('text-green-600')
  })

  // Test responsive design
  it('shows mobile view on small screens', () => {
    const wrapper = mount(StockTable, {
      props: {
        stocks: mockStocks,
        isLoading: false,
        error: null
      }
    })
    expect(wrapper.find('.md\\:hidden').exists()).toBe(true)
    expect(wrapper.find('.hidden.md\\:block').exists()).toBe(true)
  })
})