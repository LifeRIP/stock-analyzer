package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/liferip/stock-analyzer/backend/internal/models"
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
		h.logger.Error("Error getting stocks", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Error getting stocks")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items": stocks,
	})
}

// GetRecommendations maneja la solicitud para obtener recomendaciones de stocks
func (h *StockHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verificar si existe el parámetro de tiempo
	time := r.URL.Query().Get("time")

	var recommendations []models.StockRecommendation
	var err error

	if time == "" {
		// Sin parámetro de tiempo, obtener todas las recomendaciones
		recommendations, err = h.stockService.GetRecommendations(ctx)
		if err != nil {
			h.logger.Error("Error getting recommendations", zap.Error(err))
			respondWithError(w, http.StatusInternalServerError, "Error getting recommendations")
			return
		}
	} else {
		// Con parámetro de tiempo, obtener recomendaciones filtradas
		recommendations, err = h.stockService.GetRecommendationsByTime(ctx, time)
		if err != nil {
			h.logger.Error("Error getting recommendations by time", zap.Error(err))
			respondWithError(w, http.StatusInternalServerError, "Error getting recommendations by time")
			return
		}
	}

	respondWithJSON(w, http.StatusOK, recommendations)
}

// SyncStocks maneja la solicitud para sincronizar stocks desde la API externa
func (h *StockHandler) SyncStocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := h.stockService.SyncStocksFromAPI(ctx)
	if err != nil {
		h.logger.Error("Error synchronizing stocks", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Error synchronizing stocks")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Stocks synchronized correctly",
		"count":   count,
	})
}

// respondWithJSON envía una respuesta JSON al cliente
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error processing response"))
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
