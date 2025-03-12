package repository

import (
	"context"
	"errors"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/liferip/stock-analyzer/backend/internal/models"
)

// Module proporciona las dependencias del repositorio
var Module = fx.Provide(NewStockRepository)

// StockRepository interfaz que define las operaciones del repositorio
type StockRepository interface {
	GetAll(ctx context.Context) ([]models.Stock, error)
	GetAllByTime(ctx context.Context, time string) ([]models.Stock, error)
	GetByTicker(ctx context.Context, ticker string) (*models.Stock, error)
	Create(ctx context.Context, stock *models.Stock) error
	Update(ctx context.Context, stock *models.Stock) error
	Delete(ctx context.Context, id string) error
}

// stockRepository implementación de StockRepository con GORM
type stockRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewStockRepository crea una nueva instancia de StockRepository
func NewStockRepository(db *gorm.DB, logger *zap.Logger) StockRepository {
	return &stockRepository{
		db:     db,
		logger: logger.Named("stock_repository"),
	}
}

// GetAll obtiene todos los stocks de la base de datos
func (r *stockRepository) GetAll(ctx context.Context) ([]models.Stock, error) {
	var stocks []models.Stock

	result := r.db.WithContext(ctx).
		Order("time DESC").
		Find(&stocks)

	if result.Error != nil {
		r.logger.Error("Error getting stocks", zap.Error(result.Error))
		return nil, result.Error
	}

	return stocks, nil
}

// GetByTicker obtiene un stock por su ticker
func (r *stockRepository) GetByTicker(ctx context.Context, ticker string) (*models.Stock, error) {
	var stock models.Stock

	result := r.db.WithContext(ctx).
		Where("ticker = ?", ticker).
		Order("time DESC").
		First(&stock)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("Error getting stock by ticker",
			zap.String("ticker", ticker),
			zap.Error(result.Error))
		return nil, result.Error
	}

	return &stock, nil
}

// GetByTime obtiene stocks por su fecha
func (r *stockRepository) GetAllByTime(ctx context.Context, time string) ([]models.Stock, error) {
	var stocks []models.Stock

	result := r.db.WithContext(ctx).
		Where("time BETWEEN ? AND ?", time+" 00:00:00", time+" 23:59:59").
		Order("time DESC").
		Find(&stocks)

	if result.Error != nil {
		r.logger.Error("Error getting stocks by date",
			zap.String("time", time),
			zap.Error(result.Error))
		return nil, result.Error
	}

	return stocks, nil
}

// Create crea un nuevo stock en la base de datos
func (r *stockRepository) Create(ctx context.Context, stock *models.Stock) error {
	result := r.db.WithContext(ctx).Create(stock)

	if result.Error != nil {
		r.logger.Error("Error creating stock",
			zap.String("ticker", stock.Ticker),
			zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// Update actualiza un stock existente
func (r *stockRepository) Update(ctx context.Context, stock *models.Stock) error {
	result := r.db.WithContext(ctx).Save(stock)

	if result.Error != nil {
		r.logger.Error("Error updating stock",
			zap.String("id", stock.ID),
			zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// Delete elimina un stock por su ID
func (r *stockRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Delete(&models.Stock{}, "id = ?", id)

	if result.Error != nil {
		r.logger.Error("Error deleting stock",
			zap.String("id", id),
			zap.Error(result.Error))
		return result.Error
	}

	return nil
}
