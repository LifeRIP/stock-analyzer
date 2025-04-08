<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useStockStore } from "../stores/stockStore";
import { ArrowRightIcon } from "@heroicons/vue/24/outline";
import RecommendationCard from "../components/RecommendationCard.vue";

const stockStore = useStockStore();
const dateFilter = ref("");
const recommendations = computed(() => stockStore.recommendations);

onMounted(async () => {
  await stockStore.fetchRecommendations();
});

const fetchFilteredRecommendations = async () => {
  if (dateFilter.value) {
    await stockStore.fetchRecommendations(dateFilter.value);
  } else {
    await stockStore.fetchRecommendations();
  }
};

const clearDateFilter = async () => {
  dateFilter.value = "";
  await stockStore.fetchRecommendations();
};

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
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Stock Recommendations</h1>

    <div class="mb-6 flex flex-col md:flex-row gap-4">
      <div class="w-full md:w-1/3">
        <label for="date-filter" class="block text-sm font-medium text-gray-700 mb-1"
          >Filter by Date</label
        >
        <input
          id="date-filter"
          v-model="dateFilter"
          type="date"
          class="w-full px-4 py-2 border border-gray-200 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500 bg-white"
          @change="fetchFilteredRecommendations"
        />
      </div>

      <div class="flex items-end">
        <button
          @click="clearDateFilter"
          class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-green-500 cursor-pointer"
        >
          Clear Filter
        </button>
      </div>
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

    <div v-else-if="recommendations.length === 0" class="text-center my-8 text-gray-500">
      No recommendations found for the selected criteria.
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <RecommendationCard
        v-for="(recommendation, index) in recommendations"
        :key="recommendation.stock.id"
        :recommendation="recommendation"
        :index="index"
      />
    </div>
  </div>
</template>
