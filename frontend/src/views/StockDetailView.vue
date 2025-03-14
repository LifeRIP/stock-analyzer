<template>
  <div>
    <div class="mb-4">
      <router-link to="/" class="text-green-600 hover:text-green-800 flex items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z" clip-rule="evenodd" />
        </svg>
        Back to Stocks
      </router-link>
    </div>
    
    <div v-if="isLoading" class="flex justify-center my-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-green-500"></div>
    </div>
    
    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative my-4">
      {{ error }}
    </div>
    
    <div v-else-if="stock" class="bg-white rounded-lg shadow overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">{{ stock.ticker }}</h1>
            <p class="text-gray-600">{{ stock.company }}</p>
          </div>
          <div class="mt-2 md:mt-0">
            <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
              Last Updated: {{ formatDate(stock.updated_at) }}
            </span>
          </div>
        </div>
      </div>
      
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="bg-gray-50 p-4 rounded-lg">
            <h2 class="text-lg font-medium text-gray-900 mb-4">Brokerage Information</h2>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600">Brokerage:</span>
                <span class="font-medium">{{ stock.brokerage }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Action:</span>
                <span class="font-medium">{{ stock.action }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Date:</span>
                <span class="font-medium">{{ formatDate(stock.time) }}</span>
              </div>
            </div>
          </div>
          
          <div class="bg-gray-50 p-4 rounded-lg">
            <h2 class="text-lg font-medium text-gray-900 mb-4">Rating & Target</h2>
            <div class="space-y-3">
              <div class="flex justify-between" v-if="stock.rating_from && stock.rating_from !== stock.rating_to">
                <span class="text-gray-600">Rating From:</span>
                <span class="font-medium">{{ stock.rating_from }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Rating:</span>
                <span class="font-medium" :class="getRatingClass(stock.rating_from, stock.rating_to)">
                  {{ stock.rating_to || 'N/A' }}
                </span>
              </div>
              <div class="flex justify-between" v-if="stock.target_from !== stock.target_to">
                <span class="text-gray-600">Target From:</span>
                <span class="font-medium">{{ stock.target_from }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Target:</span>
                <span class="font-medium" :class="getTargetClass(stock.target_from, stock.target_to)">
                  {{ stock.target_to }}
                </span>
              </div>
              <div v-if="targetChange !== 0" class="flex justify-between">
                <span class="text-gray-600">Target Change:</span>
                <span class="font-medium" :class="targetChange > 0 ? 'text-green-600' : 'text-red-600'">
                  {{ targetChange > 0 ? '+' : '' }}{{ targetChange.toFixed(2) }}%
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="text-center my-8 text-gray-500">
      Stock not found.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useStockStore } from '../store/stockStore'
import type { Stock } from '@/types'

const props = defineProps<{
  ticker: string
}>()

const stockStore = useStockStore()
const stock = ref<Stock | null>(null)
const isLoading = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
  isLoading.value = true
  try {
    const result = await stockStore.fetchStockByTicker(props.ticker)
    stock.value = result
  } catch (err) {
    error.value = 'Failed to fetch stock details'
    console.error(err)
  } finally {
    isLoading.value = false
  }
})

const formatDate = (dateString: string) => {
  if (!dateString || dateString === '0001-01-01T00:00:00Z') return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

const getRatingClass = (from: string, to: string) => {
  if (!from || !to) return 'text-gray-600'
  
  const fromLevel = stockStore.getRatingLevel(from)
  const toLevel = stockStore.getRatingLevel(to)
  
  if (toLevel > fromLevel) return 'text-green-600'
  if (toLevel < fromLevel) return 'text-red-600'
  return 'text-gray-600'
}

const getTargetClass = (from: string, to: string) => {
  const fromValue = parseFloat(from.replace('$', ''))
  const toValue = parseFloat(to.replace('$', ''))
  
  if (toValue > fromValue) return 'text-green-600'
  if (toValue < fromValue) return 'text-red-600'
  return 'text-gray-600'
}

const targetChange = computed(() => {
  if (!stock.value) return 0
  
  const fromValue = parseFloat(stock.value.target_from.replace('$', ''))
  const toValue = parseFloat(stock.value.target_to.replace('$', ''))
  
  if (fromValue === 0) return 0
  return ((toValue - fromValue) / fromValue) * 100
})
</script>

