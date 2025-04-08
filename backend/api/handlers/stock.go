package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
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

// @Summary		Get all stocks
// @Description	Retrieves the complete list of available stocks
// @Tags			stock
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string][]models.Stock
// @Failure		404	{object}	map[string]string	"No stocks found"
// @Failure		500	{object}	map[string]string	"Error getting stocks"
// @Router			/stock [get]
func (h *StockHandler) GetStocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stocks, err := h.stockService.GetAllStocks(ctx)
	if err != nil {
		h.logger.Error("Error getting stocks", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Error getting stocks")
		return
	}

	// Si no se encontraron stocks, responder con un error
	if len(stocks) == 0 {
		respondWithError(w, http.StatusNotFound, "No stocks found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"items": stocks,
	})
}

// @Summary		Get stock by ticker
// @Description	Retrieves information for a specific stock by its ticker symbol
// @Tags			stock
// @Accept			json
// @Produce		json
// @Param			ticker	path		string	true	"Stock ticker symbol (e.g. AAPL)"
// @Success		200		{object}	map[string]models.Stock
// @Failure		404		{object}	map[string]string	"Stock not found"
// @Failure		500		{object}	map[string]string	"Error getting stock by ticker"
// @Router			/stock/ticker/{ticker} [get]
func (h *StockHandler) GetStockByTicker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ticker := mux.Vars(r)["ticker"]

	stock, err := h.stockService.GetStockByTicker(ctx, ticker)
	if err != nil {
		h.logger.Error("Error getting stock by ticker", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Error getting stock by ticker")
		return
	}

	// Si no se encontró el stock, responder con un error
	if stock == nil {
		respondWithError(w, http.StatusNotFound, "Stock not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"item": stock,
	})
}

// @Summary		Get stock recommendations
// @Description	Retrieves stock recommendations, optionally filtered by time
// @Tags			stock
// @Accept			json
// @Produce		json
// @Param			time	query		string	false	"Time filter in YY/MM/DD format (e.g. '2025-03-31')"
// @Success		200		{array}		models.StockRecommendation
// @Failure		404		{object}	map[string]string	"No recommendations found"
// @Failure		500		{object}	map[string]string	"Error getting recommendations"
// @Router			/stock/recommendations [get]
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

	// Si no se encontraron recomendaciones, responder con un error
	if len(recommendations) == 0 {
		respondWithError(w, http.StatusNotFound, "No recommendations found")
		return
	}

	respondWithJSON(w, http.StatusOK, recommendations)
}

// @Summary		Synchronize stocks
// @Description	Synchronizes stocks from an external API
// @Tags			stock
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"Success message and count"
// @Failure		500	{object}	map[string]string		"Error synchronizing stocks"
// @Router			/stock/sync [post]
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
