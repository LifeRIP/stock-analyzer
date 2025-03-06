package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"

	"github.com/liferip/stock-analyzer/backend/api"
	"github.com/liferip/stock-analyzer/backend/api/handlers"
)

// RegisterStockRoutes registra las rutas relacionadas con stocks
func RegisterStockRoutes(router *mux.Router, stockHandler *handlers.StockHandler) {
	// Rutas para stocks
	router.HandleFunc("/stocks", stockHandler.GetStocks).Methods(http.MethodGet)
	router.HandleFunc("/recommendations", stockHandler.GetRecommendations).Methods(http.MethodGet)
	router.HandleFunc("/sync", stockHandler.SyncStocks).Methods(http.MethodPost)
}

// Module proporciona las dependencias de las rutas
var Module = fx.Provide(
	func() api.RegisterRoutesFn {
		return RegisterStockRoutes
	},
)
