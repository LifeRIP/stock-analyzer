import { useStockStore } from "../stores/stockStore";

export function useStockUtils() {
  const stockStore = useStockStore();
  
  const formatDate = (dateString: string) => {
    if (!dateString || dateString === "0001-01-01T00:00:00Z") return "N/A";
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const getRatingClass = (from: string, to: string) => {
    if (!from || !to) return "text-gray-500";

    const fromLevel = stockStore.getRatingLevel(from);
    const toLevel = stockStore.getRatingLevel(to);

    if (toLevel > fromLevel) return "text-green-600";
    if (toLevel < fromLevel) return "text-red-600";
    return "text-gray-500";
  };

  const getTargetClass = (from: string, to: string) => {
    if (!from || !to) return "text-gray-500";
    const fromValue = parseFloat(from.replace(/[$,]/g, ""));
    const toValue = parseFloat(to.replace(/[$,]/g, ""));

    if (toValue > fromValue) return "text-green-600";
    if (toValue < fromValue) return "text-red-600";
    return "text-gray-500";
  };

  const calculateTargetChange = (from: string, to: string) => {
    if (!from || !to) return 0;
    const fromValue = parseFloat(from.replace(/[$,]/g, ""));
    const toValue = parseFloat(to.replace(/[$,]/g, ""));

    if (isNaN(fromValue) || isNaN(toValue) || fromValue === 0) return 0;
    return ((toValue - fromValue) / fromValue) * 100;
  };

  return {
    formatDate,
    getRatingClass,
    getTargetClass,
    calculateTargetChange
  };
}