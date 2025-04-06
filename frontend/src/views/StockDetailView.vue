<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useStockStore } from "../stores/stockStore";
import { useStockUtils } from "../composables/useStockUtils";
import type { Stock } from "@/types";
import { ArrowLeftIcon } from "@heroicons/vue/24/outline";

const props = defineProps<{
  ticker: string;
}>();

const stockStore = useStockStore();
const { formatDate, getRatingClass, getTargetClass, calculateTargetChange } = useStockUtils();
const stock = ref<Stock | null>(null);

onMounted(async () => {
    const result = await stockStore.fetchStockByTicker(props.ticker);
    stock.value = result;
});

const targetChange = computed(() => {
  if (!stock.value) return 0;
  return calculateTargetChange(stock.value.target_from, stock.value.target_to);
});
</script>

<template>
  <div>
    <div class="mb-4">
      <router-link to="/" class="text-green-600 hover:text-green-800 flex items-center">
        <ArrowLeftIcon class="h-5 w-5 mr-1" />
        Back to Stocks
      </router-link>
    </div>

    <div v-if="stockStore.isLoading" class="flex justify-center my-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-green-500"></div>
    </div>

    <div
      v-else-if="stockStore.error"
      class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative my-4"
    >
      {{ stockStore.error }}
    </div>

    <div v-else-if="stock" class="bg-white rounded-lg shadow overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">{{ stock.ticker }}</h1>
            <p class="text-gray-600">{{ stock.company }}</p>
          </div>
          <div class="mt-2 md:mt-0">
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800"
            >
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
              <div
                class="flex justify-between"
                v-if="stock.rating_from && stock.rating_from !== stock.rating_to"
              >
                <span class="text-gray-600">Rating From:</span>
                <span class="font-medium">{{ stock.rating_from }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Rating:</span>
                <span
                  class="font-medium"
                  :class="getRatingClass(stock.rating_from, stock.rating_to)"
                >
                  {{ stock.rating_to || "N/A" }}
                </span>
              </div>
              <div class="flex justify-between" v-if="stock.target_from !== stock.target_to">
                <span class="text-gray-600">Target From:</span>
                <span class="font-medium">{{ stock.target_from }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Target:</span>
                <span
                  class="font-medium"
                  :class="getTargetClass(stock.target_from, stock.target_to)"
                >
                  {{ stock.target_to }}
                </span>
              </div>
              <div v-if="targetChange !== 0" class="flex justify-between">
                <span class="text-gray-600">Target Change:</span>
                <span
                  class="font-medium"
                  :class="targetChange > 0 ? 'text-green-600' : 'text-red-600'"
                >
                  {{ targetChange > 0 ? "+" : "" }}{{ targetChange.toFixed(2) }}%
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-center my-8 text-gray-500">Stock not found.</div>
  </div>
</template>
