<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Stock Recommendations</h1>
    
    <div class="mb-6 flex flex-col md:flex-row gap-4">
      <div class="w-full md:w-1/3">
        <label for="date-filter" class="block text-sm font-medium text-gray-700 mb-1">Filter by Date</label>
        <input 
          id="date-filter"
          v-model="dateFilter" 
          type="date" 
          class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-green-500"
          @change="fetchFilteredRecommendations"
        />
      </div>
      
      <div class="flex items-end">
        <button 
          @click="clearDateFilter" 
          class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-green-500"
        >
          Clear Filter
        </button>
      </div>
    </div>
    
    <div v-if="stockStore.isLoading" class="flex justify-center my-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-green-500"></div>
    </div>
    
    <div v-else-if="stockStore.error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative my-4">
      {{ stockStore.error }}
    </div>
    
    <div v-else-if="recommendations.error" class="text-center my-8 text-gray-500">
      No recommendations found for the selected criteria.
    </div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="(recommendation, index) in recommendations" 
        :key="recommendation.stock.id"
        class="bg-white rounded-lg shadow overflow-hidden hover:shadow-lg transition-shadow duration-300"
      >
        <div class="p-4 bg-gradient-to-r from-green-500 to-green-600 text-white">
          <div class="flex justify-between items-center">
            <h2 class="text-xl font-bold">{{ recommendation.stock.ticker }}</h2>
            <span class="text-sm bg-white text-green-600 px-2 py-1 rounded-full font-bold">
              Rank #{{ index + 1 }}
            </span>
          </div>
          <p class="text-green-100 text-sm truncate">{{ recommendation.stock.company }}</p>
        </div>
        
        <div class="p-4">
          <div class="mb-4">
            <div class="flex justify-between mb-1">
              <span class="text-sm font-medium text-gray-700">Score</span>
              <span class="text-sm font-bold">{{ recommendation.score.toFixed(2) }}</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div class="bg-green-600 h-2 rounded-full" :style="{ width: `${Math.min(recommendation.score, 100)}%` }"></div>
            </div>
          </div>
          
          <div class="mb-4">
            <div class="flex justify-between mb-1">
              <span class="text-sm font-medium text-gray-700">Potential Upside</span>
              <span class="text-sm font-bold text-green-600">
                {{ recommendation.potential_up != null ? '+' + recommendation.potential_up.toFixed(2) : '0.00' }}%
              </span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div class="bg-green-600 h-2 rounded-full" 
                   :style="{ width: `${recommendation.potential_up != null ? Math.min(recommendation.potential_up, 100) : 0}%` }">
              </div>
            </div>
          </div>
          
          <div class="mb-4">
            <h3 class="text-sm font-medium text-gray-700 mb-2">Why We Recommend</h3>
            <ul class="text-sm text-gray-600 space-y-1 pl-5 list-disc">
              <li v-for="(reason, i) in recommendation.reasons" :key="i">
                {{ reason }}
              </li>
            </ul>
          </div>
          
          <div class="mt-4 pt-4 border-t border-gray-200">
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">{{ recommendation.stock.brokerage }}</span>
              <span class="text-gray-500">{{ formatDate(recommendation.stock.time) }}</span>
            </div>
            <div class="flex justify-between text-sm mt-2">
              <span class="text-gray-500">
                Target: 
                <template v-if="recommendation.stock.target_from !== recommendation.stock.target_to">
                  {{ recommendation.stock.target_from }} → 
                </template>
                <span class="text-green-600 font-medium">{{ recommendation.stock.target_to }}</span>
              </span>
            </div>
            <div class="flex justify-between text-sm mt-2" v-if="recommendation.stock.rating_from || recommendation.stock.rating_to">
              <span class="text-gray-500">
                Rating: 
                <template v-if="recommendation.stock.rating_from && recommendation.stock.rating_from !== recommendation.stock.rating_to">
                  {{ recommendation.stock.rating_from }} → 
                </template>
                <span :class="getRatingClass(recommendation.stock.rating_from, recommendation.stock.rating_to)" class="font-medium">
                  {{ recommendation.stock.rating_to || 'N/A' }}
                </span>
              </span>
            </div>
          </div>
        </div>
        
        <div class="px-4 py-3 bg-gray-50 border-t border-gray-200">
          <router-link 
            :to="{ name: 'stock-detail', params: { ticker: recommendation.stock.ticker }}" 
            class="text-green-600 hover:text-green-800 text-sm font-medium flex items-center justify-center"
          >
            View Details
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useStockStore } from '../store/stockStore'
import type { StockRecommendation } from '@/types'

const stockStore = useStockStore()
const dateFilter = ref('')
const recommendations = computed(() => stockStore.recommendations)

onMounted(async () => {
  await stockStore.fetchRecommendations()
})

const fetchFilteredRecommendations = async () => {
  if (dateFilter.value) {
    await stockStore.fetchRecommendations(dateFilter.value)
  } else {
    await stockStore.fetchRecommendations()
  }
}

const clearDateFilter = async () => {
  dateFilter.value = ''
  await stockStore.fetchRecommendations()
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

// Add the getRatingClass function
const getRatingClass = (from: string, to: string) => {
  if (!from || !to) return 'text-gray-600'
  
  const fromLevel = stockStore.getRatingLevel(from)
  const toLevel = stockStore.getRatingLevel(to)
  
  if (toLevel > fromLevel) return 'text-green-600'
  if (toLevel < fromLevel) return 'text-red-600'
  return 'text-gray-600'
}
</script>

