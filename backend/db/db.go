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

	// Primero intentamos conectarnos a la base de datos postgres (default)
	tempConfig := *cfg
	tempConfig.DatabaseName = "defaultdb"

	tempDB, err := gorm.Open(postgres.Open(createURL(&tempConfig)), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres database: %w", err)
	}

	// Crear la base de datos si no existe
	createDBQuery := fmt.Sprintf(
		"SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')",
		cfg.DatabaseName, cfg.DatabaseName,
	)

	// Ejecutar la consulta para crear la base de datos
	var result string
	tempDB.Raw(createDBQuery).Scan(&result)

	fmt.Println("result", result)
	if result != "" {
		// Si la base de datos no existe, crearla
		tempDB.Exec(result)
		log.Printf("Database %s created successfully", cfg.DatabaseName)
	}

	// Cerrar la conexión temporal
	sqlDB, _ := tempDB.DB()
	sqlDB.Close()

	// Conectar a la base de datos destino
	db, err := gorm.Open(postgres.Open(createURL(cfg)), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Configurar el pool de conexiones
	sqlDB, err = db.DB()
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

// createUrl crea la URL de conexión a la base de datos
func createURL(cfg *config.Config) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DatabaseUser,
		cfg.DatabasePass,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
		cfg.DatabaseSSLMode,
	)
}

// migrateSchema migra el esquema de la base de datos
func migrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Stock{},
	)
}
