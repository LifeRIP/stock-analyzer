<script setup lang="ts">
import { computed } from "vue";
import { ArrowRightIcon } from "@heroicons/vue/24/outline";
import { useStockStore } from "../stores/stockStore";

const props = defineProps({
  recommendation: {
    type: Object,
    required: true
  },
  index: {
    type: Number,
    required: true
  }
});

const stockStore = useStockStore();

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
  <div
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
          <div
            class="bg-green-600 h-2 rounded-full"
            :style="{ width: `${Math.min(recommendation.score, 100)}%` }"
          ></div>
        </div>
      </div>

      <div class="mb-4">
        <div class="flex justify-between mb-1">
          <span class="text-sm font-medium text-gray-700">Potential Upside</span>
          <span class="text-sm font-bold text-green-600">
            {{
              recommendation.potential_up != null
                ? "+" + recommendation.potential_up.toFixed(2)
                : "0.00"
            }}%
          </span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div
            class="bg-green-600 h-2 rounded-full"
            :style="{
              width: `${
                recommendation.potential_up != null
                  ? Math.min(recommendation.potential_up, 100)
                  : 0
              }%`,
            }"
          ></div>
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
            <template
              v-if="recommendation.stock.target_from !== recommendation.stock.target_to"
            >
              {{ recommendation.stock.target_from }} →
            </template>
            <span class="text-green-600 font-medium">{{ recommendation.stock.target_to }}</span>
          </span>
        </div>
        <div
          class="flex justify-between text-sm mt-2"
          v-if="recommendation.stock.rating_from || recommendation.stock.rating_to"
        >
          <span class="text-gray-500">
            Rating:
            <template
              v-if="
                recommendation.stock.rating_from &&
                recommendation.stock.rating_from !== recommendation.stock.rating_to
              "
            >
              {{ recommendation.stock.rating_from }} →
            </template>
            <span
              :class="
                getRatingClass(recommendation.stock.rating_from, recommendation.stock.rating_to)
              "
              class="font-medium"
            >
              {{ recommendation.stock.rating_to || "N/A" }}
            </span>
          </span>
        </div>
      </div>
    </div>

    <div class="px-4 py-3 bg-gray-50 border-t border-gray-200">
      <router-link
        :to="{ name: 'stock-detail', params: { ticker: recommendation.stock.ticker } }"
        class="text-green-600 hover:text-green-800 text-sm font-medium flex items-center justify-center"
      >
        View Details
        <ArrowRightIcon class="h-4 w-4 ml-1" />
      </router-link>
    </div>
  </div>
</template>