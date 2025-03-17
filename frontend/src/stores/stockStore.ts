import { defineStore } from "pinia";
import { ref } from "vue";
import type { Stock, StockRecommendation } from "@/types";

export const useStockStore = defineStore("stock", () => {
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8081";
  const stocks = ref<Stock[]>([]);
  const recommendations = ref<StockRecommendation[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const fetchStocks = async () => {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await fetch(`${API_URL}/api/stock`);
      const data = await response.json();
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      stocks.value = data.items;
      //console.log(data)
      return data;
    } catch (err) {
      error.value = "Failed to fetch stocks";
      console.error(err);
    } finally {
      isLoading.value = false;
    }
  };

  const fetchStockByTicker = async (ticker: string) => {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await fetch(`${API_URL}/api/stock/ticker/${ticker}`);
      const data = await response.json();
      return data.item;
    } catch (err) {
      error.value = "Failed to fetch stock details";
      console.error(err);
      return null;
    } finally {
      isLoading.value = false;
    }
  };

  const fetchRecommendations = async (date?: string) => {
    isLoading.value = true;
    error.value = null;
    try {
      const url = date
        ? `${API_URL}/api/stock/recommendations?time=${date}`
        : `${API_URL}/api/stock/recommendations`;

      const response = await fetch(url);
      const data = await response.json();
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      recommendations.value = data;
      return data;
    } catch (err) {
      error.value = "Failed to fetch recommendations";
      console.error(err);
      return [];
    } finally {
      isLoading.value = false;
    }
  };

  const syncStocks = async () => {
    isLoading.value = true;
    error.value = null;
    try {
      await fetch(`${API_URL}/api/stock/sync`, {
        method: "POST",
      });
      await fetchStocks();
      return true;
    } catch (err) {
      error.value = "Failed to sync stocks";
      console.error(err);
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  // Agregar la función para obtener el nivel de rating
  const getRatingLevel = (rating: string): number => {
    const positiveRatings = [
      "Strong-Buy",
      "Buy",
      "Speculative Buy",
      "Positive",
      "Market Outperform",
      "Sector Outperform",
      "Outperform",
      "Outperformer",
      "Overweight",
    ];
    const neutralRatings = [
      "Neutral",
      "Hold",
      "Equal Weight",
      "In-Line",
      "Inline",
      "Sector Perform",
      "Market Perform",
      "Sector Weight",
    ];
    const negativeRatings = [
      "Underweight",
      "Underperform",
      "Sector Underperform",
      "Reduce",
      "Negative",
    ];

    if (positiveRatings.includes(rating)) return 3;
    if (neutralRatings.includes(rating)) return 2;
    if (negativeRatings.includes(rating)) return 1;
    return 0;
  };

  // Agregar la función al return del store
  return {
    stocks,
    recommendations,
    isLoading,
    error,
    fetchStocks,
    fetchStockByTicker,
    fetchRecommendations,
    syncStocks,
    getRatingLevel,
  };
});
