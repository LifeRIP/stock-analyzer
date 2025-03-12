package logger

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/liferip/stock-analyzer/backend/config"
)

// Module proporciona las dependencias del logger
var Module = fx.Provide(NewLogger)

// NewLogger crea una nueva instancia del logger
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var config zap.Config

	if cfg.Environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Configurar salida
	config.OutputPaths = []string{"stdout"}
	if cfg.Environment == "production" {
		config.OutputPaths = append(config.OutputPaths, "logs/app.log")
	}

	// Crear logger
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	// Reemplazar el logger global de zap
	zap.ReplaceGlobals(logger)

	return logger, nil
}

// Sync sincroniza el logger antes de que la aplicación se cierre
func Sync(logger *zap.Logger) {
	if err := logger.Sync(); err != nil {
		// Ignorar errores de sincronización en stdout/stderr
		if err.Error() != "sync /dev/stdout: invalid argument" {
			os.Stderr.WriteString("Error synchronizing logger: " + err.Error() + "\n")
		}
	}
}
