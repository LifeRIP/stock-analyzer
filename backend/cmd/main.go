package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/liferip/stock-analyzer/backend/api"
	"github.com/liferip/stock-analyzer/backend/api/handlers"
	"github.com/liferip/stock-analyzer/backend/api/routes"
	"github.com/liferip/stock-analyzer/backend/config"
	"github.com/liferip/stock-analyzer/backend/db"
	"github.com/liferip/stock-analyzer/backend/internal/repository"
	"github.com/liferip/stock-analyzer/backend/internal/service"
	"github.com/liferip/stock-analyzer/backend/pkg/httpclient"
	"github.com/liferip/stock-analyzer/backend/pkg/logger"
)

func main() {
	fx.New(
		// Incluir módulos
		config.Module,
		logger.Module,
		db.Module,
		httpclient.Module,
		repository.Module,
		service.Module,
		handlers.Module,
		routes.Module,
		api.Module,

		// Registrar hooks
		fx.Invoke(register),
	).Run()
}

func register(
	lc fx.Lifecycle,
	cfg *config.Config,
	router *mux.Router,
	logger *zap.Logger,
	stockService service.StockService,
) {
	// Crear servidor HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Registrar hooks del ciclo de vida
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Iniciar servidor en una goroutine
			go func() {
				logger.Info("Servidor iniciado", zap.String("port", cfg.ServerPort))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Error al iniciar el servidor", zap.Error(err))
				}
			}()

			// Sincronizar stocks al inicio (opcional)
			go func() {
				logger.Info("Sincronizando stocks desde la API...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				defer cancel()

				count, err := stockService.SyncStocksFromAPI(ctx)
				if err != nil {
					logger.Error("Error al sincronizar stocks", zap.Error(err))
				} else {
					logger.Info("Sincronización completada", zap.Int("count", count))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Crear un contexto con timeout para el apagado
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Cerrar el servidor HTTP
			logger.Info("Cerrando servidor HTTP...")
			if err := server.Shutdown(shutdownCtx); err != nil {
				logger.Error("Error al cerrar el servidor", zap.Error(err))
				return err
			}

			// Sincronizar logger
			logger.Sync()
			return nil
		},
	})

	// Manejar señales de cierre
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		logger.Info("Señal recibida, iniciando apagado", zap.String("signal", sig.String()))
		os.Exit(0)
	}()
}
