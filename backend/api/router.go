package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/liferip/stock-analyzer/backend/api/handlers"
)

// RegisterRoutesFn es un tipo para funciones que registran rutas
type RegisterRoutesFn func(*mux.Router, *handlers.StockHandler)

// NewRouter crea un nuevo router con todas las rutas configuradas
func NewRouter(
	logger *zap.Logger,
	stockHandler *handlers.StockHandler,
	registerStockRoutes RegisterRoutesFn,
) *mux.Router {
	router := mux.NewRouter()

	// Middleware healthcheck
	router.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Middleware para CORS
	router.Use(corsMiddleware)

	// Middleware para logging
	router.Use(loggingMiddleware(logger))

	// API endpoints
	api := router.PathPrefix("/api").Subrouter()

	// Registrar rutas
	registerStockRoutes(api, stockHandler)

	return router
}

// corsMiddleware agrega los headers CORS necesarios
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware registra información sobre cada solicitud HTTP
func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote", r.RemoteAddr),
			)
			next.ServeHTTP(w, r)
		})
	}
}

// Module proporciona las dependencias del router
var Module = fx.Provide(NewRouter)
