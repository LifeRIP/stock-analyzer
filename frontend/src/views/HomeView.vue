<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Stock List</h1>
    
    <StockTable 
      :stocks="paginatedStocks" 
      :isLoading="stockStore.isLoading" 
      :error="stockStore.error"
      @search="handleSearch"
      @page-size-change="handlePageSizeChange"
      @stock-click="goToStockDetail"
    />
    
    <Pagination 
      v-if="filteredStocks.length > 0"
      :current-page="currentPage" 
      :page-size="pageSize" 
      :total-items="filteredStocks.length"
      @page-change="handlePageChange"
      class="mt-4"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStockStore } from '../store/stockStore'
import StockTable from '../components/StockTable.vue'
import Pagination from '../components/Pagination.vue'

const router = useRouter()
const stockStore = useStockStore()
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(25)

onMounted(async () => {
  if (stockStore.stocks.length === 0) {
    await stockStore.fetchStocks()
  }
})

const filteredStocks = computed(() => {
  if (!searchQuery.value) return stockStore.stocks
  
  const query = searchQuery.value.toLowerCase()
  return stockStore.stocks.filter(stock => 
    stock.ticker.toLowerCase().includes(query) || 
    stock.company.toLowerCase().includes(query)
  )
})

const paginatedStocks = computed(() => {
  const startIndex = (currentPage.value - 1) * pageSize.value
  return filteredStocks.value.slice(startIndex, startIndex + pageSize.value)
})

const handleSearch = (query: string) => {
  searchQuery.value = query
  currentPage.value = 1 // Reset to first page when searching
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1 // Reset to first page when changing page size
}

const handlePageChange = (page: number) => {
  currentPage.value = page
}

const goToStockDetail = (ticker: string) => {
  router.push({ name: 'stock-detail', params: { ticker } })
}
</script>

