<script setup lang="ts">
import { ref } from "vue";
import { useStockStore } from "../stores/stockStore";
import type { Stock } from "@/types";

defineProps<{
  stocks: Stock[];
  isLoading: boolean;
  error: string | null;
}>();

defineEmits<{
  (e: "search", query: string): void;
  (e: "page-size-change", size: number): void;
  (e: "stock-click", ticker: string): void;
}>();

const stockStore = useStockStore();
const searchQuery = ref("");
const pageSize = ref(25);

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString();
};

const getRatingClass = (from: string, to: string) => {
  if (!from || !to) return "text-gray-600";

  const fromLevel = stockStore.getRatingLevel(from);
  const toLevel = stockStore.getRatingLevel(to);

  if (toLevel > fromLevel) return "text-green-600";
  if (toLevel < fromLevel) return "text-red-600";
  return "text-gray-600";
};

const getTargetClass = (from: string, to: string) => {
  const fromValue = parseFloat(from.replace("$", ""));
  const toValue = parseFloat(to.replace("$", ""));

  if (toValue > fromValue) return "text-green-600";
  if (toValue < fromValue) return "text-red-600";
  return "text-gray-600";
};
</script>

<template>
  <div>
    <div class="flex flex-col md:flex-row justify-between items-center mb-4 gap-4">
      <div class="w-full md:w-1/3">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by ticker or company..."
          class="w-full px-4 py-2 border border-gray-200 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500 placeholder-gray-400 bg-white"
          @input="$emit('search', searchQuery)"
        />
      </div>

      <div class="flex items-center space-x-2">
        <label for="page-size" class="text-sm text-gray-600">Show:</label>
        <select
          id="page-size"
          v-model="pageSize"
          class="border rounded-md px-2 py-1 text-sm border border-gray-200 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500 placeholder-gray-400 bg-white"
          @change="$emit('page-size-change', pageSize)"
        >
          <option :value="10">10</option>
          <option :value="25">25</option>
          <option :value="50">50</option>
          <option :value="100">100</option>
        </select>
      </div>
    </div>

    <div v-if="isLoading" class="flex justify-center my-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-green-500"></div>
    </div>

    <div
      v-else-if="error"
      class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative my-4"
    >
      {{ error }}
    </div>

    <div v-else-if="stocks.length === 0" class="text-center my-8 text-gray-500">
      No stocks found. Try a different search or sync stocks.
    </div>

    <div v-else>
      <!-- Mobile view (under 768px) -->
      <div class="md:hidden space-y-4">
        <div
          v-for="stock in stocks"
          :key="stock.id"
          class="bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow duration-200"
          @click="$emit('stock-click', stock.ticker)"
        >
          <div class="flex justify-between items-center mb-2">
            <span class="font-medium text-green-600">{{ stock.ticker }}</span>
            <span class="text-xs text-gray-500">{{ formatDate(stock.time) }}</span>
          </div>
          <div class="text-sm mb-2 truncate">{{ stock.company }}</div>
          <div class="grid grid-cols-2 gap-2 text-xs">
            <div>
              <span class="text-gray-500">Action:</span>
              <span class="ml-1">{{ stock.action }}</span>
            </div>
            <div>
              <span class="text-gray-500">Broker:</span>
              <span class="ml-1 truncate block">{{ stock.brokerage }}</span>
            </div>
            <div>
              <span class="text-gray-500">Rating:</span>
              <span class="ml-1" :class="getRatingClass(stock.rating_from, stock.rating_to)">
                <template v-if="stock.rating_from && stock.rating_from !== stock.rating_to">
                  {{ stock.rating_from }} → {{ stock.rating_to }}
                </template>
                <template v-else>
                  {{ stock.rating_to || "N/A" }}
                </template>
              </span>
            </div>
            <div>
              <span class="text-gray-500">Target:</span>
              <span class="ml-1" :class="getTargetClass(stock.target_from, stock.target_to)">
                <template v-if="stock.target_from !== stock.target_to">
                  {{ stock.target_from }} → {{ stock.target_to }}
                </template>
                <template v-else>
                  {{ stock.target_to }}
                </template>
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Tablet and desktop view (768px and above) -->
      <div class="hidden md:block bg-white rounded-lg shadow">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 table-fixed">
            <thead class="bg-gray-50">
              <tr>
                <th
                  scope="col"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[10%]"
                >
                  Ticker
                </th>
                <th
                  scope="col"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[20%]"
                >
                  Company
                </th>
                <th
                  scope="col"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[20%]"
                >
                  Broker
                </th>
                <th
                  scope="col"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[15%]"
                >
                  Action
                </th>
                <th
                  scope="col"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[15%]"
                >
                  Rating
                </th>
                <th
                  scope="col"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[12%]"
                >
                  Target
                </th>
                <th
                  scope="col"
                  class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[8%]"
                >
                  Date
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="stock in stocks"
                :key="stock.id"
                class="hover:bg-gray-50 cursor-pointer"
                @click="$emit('stock-click', stock.ticker)"
              >
                <td class="px-3 py-3 whitespace-nowrap">
                  <div class="font-medium text-green-600">{{ stock.ticker }}</div>
                </td>
                <td class="px-3 py-3">
                  <div class="text-sm text-gray-900 truncate" :title="stock.company">
                    {{ stock.company }}
                  </div>
                </td>
                <td class="px-3 py-3">
                  <div class="text-sm text-gray-500 truncate" :title="stock.brokerage">
                    {{ stock.brokerage }}
                  </div>
                </td>
                <td class="px-3 py-3">
                  <div class="text-sm text-gray-500 truncate" :title="stock.action">
                    {{ stock.action }}
                  </div>
                </td>
                <td class="px-3 py-3">
                  <div class="text-sm" :class="getRatingClass(stock.rating_from, stock.rating_to)">
                    <template v-if="stock.rating_from && stock.rating_from !== stock.rating_to">
                      <span class="text-gray-500">{{ stock.rating_from }}</span> →
                      {{ stock.rating_to }}
                    </template>
                    <template v-else>
                      {{ stock.rating_to || "N/A" }}
                    </template>
                  </div>
                </td>
                <td class="px-3 py-3">
                  <div class="text-sm" :class="getTargetClass(stock.target_from, stock.target_to)">
                    <template v-if="stock.target_from !== stock.target_to">
                      <span class="text-gray-500">{{ stock.target_from }}</span> →
                      {{ stock.target_to }}
                    </template>
                    <template v-else>
                      {{ stock.target_to }}
                    </template>
                  </div>
                </td>
                <td class="px-3 py-3 whitespace-nowrap text-xs text-gray-500">
                  {{ formatDate(stock.time) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
