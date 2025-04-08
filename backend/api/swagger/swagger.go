package swagger

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/liferip/stock-analyzer/backend/config"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RegisterSwaggerRoutes configura las rutas para Swagger UI
func RegisterSwaggerRoutes(router *mux.Router, logger *zap.Logger, cfg *config.Config) {
	logger.Info("Registering Swagger routes")

	// Redirigir la ruta raíz a la documentación Swagger
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	}).Methods(http.MethodGet)

	// Configurar ruta para la documentación Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

// Module proporciona las dependencias para Swagger
var Module = fx.Options(
	fx.Provide(),
	fx.Invoke(RegisterSwaggerRoutes),
)
