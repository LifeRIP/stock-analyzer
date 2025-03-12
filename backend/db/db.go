package db

import (
	"fmt"
	"log"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/liferip/stock-analyzer/backend/config"
	"github.com/liferip/stock-analyzer/backend/internal/models"
)

// Module proporciona las dependencias de la base de datos
var Module = fx.Provide(NewDatabase)

// NewDatabase inicializa y configura la conexión a la base de datos
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Configurar el logger de GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Si estamos en producción, reducir el nivel de log
	if cfg.Environment == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Configurar el pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting SQL connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Migrar el esquema
	if err := migrateSchema(db); err != nil {
		return nil, fmt.Errorf("error migrating schema: %w", err)
	}

	log.Println("Connection to the database established successfully")
	return db, nil
}

// migrateSchema migra el esquema de la base de datos
func migrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Stock{},
	)
}
