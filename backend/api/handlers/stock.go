package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/liferip/stock-analyzer/backend/internal/service"
)

// Module proporciona las dependencias de los handlers
var Module = fx.Provide(NewStockHandler)

// StockHandler maneja las solicitudes relacionadas con stocks
type StockHandler struct {
	stockService service.StockService
	logger       *zap.Logger
}

// NewStockHandler crea una nueva instancia de StockHandler
func NewStockHandler(
	stockService service.StockService,
	logger *zap.Logger,
) *StockHandler {
	return &StockHandler{
		stockService: stockService,
		logger:       logger.Named("stock_handler"),
	}
}

// GetStocks maneja la solicitud para obtener todos los stocks
func (h *StockHandler) GetStocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stocks, err := h.stockService.GetAllStocks(ctx)
	if err != nil {
		h.logger.Error("Error al obtener stocks", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Error al obtener stocks")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items": stocks,
	})
}

// GetRecommendations maneja la solicitud para obtener recomendaciones de stocks
func (h *StockHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	recommendations, err := h.stockService.GetRecommendations(ctx)
	if err != nil {
		h.logger.Error("Error al obtener recomendaciones", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Error al obtener recomendaciones")
		return
	}

	respondWithJSON(w, http.StatusOK, recommendations)
}

// SyncStocks maneja la solicitud para sincronizar stocks desde la API externa
func (h *StockHandler) SyncStocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := h.stockService.SyncStocksFromAPI(ctx)
	if err != nil {
		h.logger.Error("Error al sincronizar stocks", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Error al sincronizar stocks")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Stocks sincronizados correctamente",
		"count":   count,
	})
}

// respondWithJSON envía una respuesta JSON al cliente
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al procesar la respuesta"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// respondWithError envía un error al cliente
func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
