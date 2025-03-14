package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/liferip/stock-analyzer/backend/internal/models"
	"github.com/liferip/stock-analyzer/backend/internal/repository"
	"github.com/liferip/stock-analyzer/backend/pkg/httpclient"
)

// Module proporciona las dependencias del servicio
var Module = fx.Provide(NewStockService)

// StockService interfaz que define las operaciones del servicio
type StockService interface {
	GetAllStocks(ctx context.Context) ([]models.Stock, error)
	GetStockByTicker(ctx context.Context, ticker string) (*models.Stock, error)
	SyncStocksFromAPI(ctx context.Context) (int, error)
	GetRecommendations(ctx context.Context) ([]models.StockRecommendation, error)
	GetRecommendationsByTime(ctx context.Context, time string) ([]models.StockRecommendation, error)
}

// stockService implementación de StockService
type stockService struct {
	repo        repository.StockRepository
	stockClient *httpclient.StockClient
	logger      *zap.Logger
}

// NewStockService crea una nueva instancia de StockService
func NewStockService(
	repo repository.StockRepository,
	stockClient *httpclient.StockClient,
	logger *zap.Logger,
) StockService {
	return &stockService{
		repo:        repo,
		stockClient: stockClient,
		logger:      logger.Named("stock_service"),
	}
}

// GetAllStocks obtiene todos los stocks
func (s *stockService) GetAllStocks(ctx context.Context) ([]models.Stock, error) {
	return s.repo.GetAll(ctx)
}

// GetStockByTicker obtiene un stock por su ticker
func (s *stockService) GetStockByTicker(ctx context.Context, ticker string) (*models.Stock, error) {
	return s.repo.GetByTicker(ctx, ticker)
}

// SyncStocksFromAPI sincroniza los stocks desde la API externa
func (s *stockService) SyncStocksFromAPI(ctx context.Context) (int, error) {
	timeStart := time.Now()
	var nextPage string
	var count int
	var mu sync.Mutex // Mutex para proteger el contador

	// Número de workers para procesar stocks en paralelo
	const numWorkers = 10

	for {
		// Obtener datos de la API
		response, err := s.stockClient.GetStocks(nextPage)
		if err != nil {
			s.logger.Error("Error getting stocks from the API", zap.Error(err))
			return count, fmt.Errorf("error getting stocks from the API: %w", err)
		}

		// Crear un canal para los items y waitgroup para esperar a que terminen
		itemCh := make(chan models.StockItem, len(response.Items))
		var wg sync.WaitGroup
		errCh := make(chan error, len(response.Items))

		// Iniciar workers
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for item := range itemCh {
					if err := s.processStockItem(ctx, item, &mu, &count); err != nil {
						select {
						case errCh <- err:
							// Enviar el error al canal
						default:
							// Si el canal está lleno, logear el error y continuar
							s.logger.Error("Error processing stock item", zap.Error(err))
						}
						return
					}
				}
			}()
		}

		// Enviar items a los workers
		for _, item := range response.Items {
			itemCh <- item
		}
		close(itemCh)

		// Esperar a que terminen todos los workers
		wg.Wait()

		// Verificar si hay errores
		select {
		case err := <-errCh:
			s.logger.Error("Error in worker", zap.Error(err))
			return count, err
		default:
			// No hay errores, continuamos
		}

		// Verificar si hay más páginas
		nextPage = response.NextPage
		if nextPage == "" {
			break
		}

		// Añadir delay para no saturar la API
		//time.Sleep(500 * time.Millisecond)
	}

	s.logger.Info("Synchronization completed", zap.Int("count", count), zap.Duration("duration", time.Since(timeStart)))
	return count, nil
}

// processStockItem procesa un solo item de stock
func (s *stockService) processStockItem(ctx context.Context, item models.StockItem, mu *sync.Mutex, count *int) error {
	// Parsear la fecha
	timeValue, err := time.Parse(time.RFC3339, item.Time)
	if err != nil {
		s.logger.Warn("Error parsing date, using current time",
			zap.String("time", item.Time),
			zap.Error(err))
		timeValue = time.Now()
	}

	// Crear modelo de stock
	stock := &models.Stock{
		Ticker:     item.Ticker,
		Company:    item.Company,
		Brokerage:  item.Brokerage,
		Action:     item.Action,
		RatingFrom: item.RatingFrom,
		RatingTo:   item.RatingTo,
		TargetFrom: item.TargetFrom,
		TargetTo:   item.TargetTo,
		Time:       timeValue,
	}

	// Verificar si ya existe
	existing, err := s.repo.GetByTickerSimple(ctx, item.Ticker)
	if err != nil {
		s.logger.Error("Error checking existing stock",
			zap.String("ticker", item.Ticker),
			zap.Error(err))
		return fmt.Errorf("error checking existing stock: %w", err)
	}

	if existing == nil {
		// Crear nuevo stock
		if err := s.repo.Create(ctx, stock); err != nil {
			s.logger.Error("Error creating stock",
				zap.String("ticker", item.Ticker),
				zap.Error(err))
			return fmt.Errorf("error creating stock: %w", err)
		}
	} else {
		// Truncar la fecha a segundos para evitar problemas de precisión
		timeValue = timeValue.Truncate(time.Second)
		existing.Time = existing.Time.Truncate(time.Second)

		// Actualizar stock existente si la fecha es más reciente
		if timeValue.After(existing.Time) {
			stock.ID = existing.ID
			if err := s.repo.Update(ctx, stock); err != nil {
				s.logger.Error("Error updating stock:",
					zap.String("ticker", item.Ticker),
					zap.Error(err))
				return fmt.Errorf("error updating stock: %w", err)
			}
		}
	}

	mu.Lock()
	*count++
	mu.Unlock()

	return nil
}

// // SyncStocksFromAPI sincroniza los stocks desde la API externa
// func (s *stockService) SyncStocksFromAPI(ctx context.Context) (int, error) {
// 	timeStart := time.Now()
// 	var nextPage string
// 	var count int

// 	for {
// 		// Obtener datos de la API
// 		response, err := s.stockClient.GetStocks(nextPage)
// 		if err != nil {
// 			s.logger.Error("Error getting stocks from the API", zap.Error(err))
// 			return count, fmt.Errorf("error getting stocks from the API: %w", err)
// 		}

// 		// Guardar cada stock en la base de datos
// 		for _, item := range response.Items {

// 			// Parsear la fecha
// 			timeValue, err := time.Parse(time.RFC3339, item.Time)

// 			if err != nil {
// 				s.logger.Warn("Error parsing date, using current time",
// 					zap.String("time", item.Time),
// 					zap.Error(err))
// 				timeValue = time.Now()
// 			}

// 			// Crear modelo de stock
// 			stock := &models.Stock{
// 				Ticker:     item.Ticker,
// 				Company:    item.Company,
// 				Brokerage:  item.Brokerage,
// 				Action:     item.Action,
// 				RatingFrom: item.RatingFrom,
// 				RatingTo:   item.RatingTo,
// 				TargetFrom: item.TargetFrom,
// 				TargetTo:   item.TargetTo,
// 				Time:       timeValue,
// 			}

// 			// Verificar si ya existe
// 			existing, err := s.repo.GetByTickerSimple(ctx, item.Ticker)
// 			if err != nil {
// 				s.logger.Error("Error checking existing stock",
// 					zap.String("ticker", item.Ticker),
// 					zap.Error(err))
// 				return count, fmt.Errorf("error checking existing stock: %w", err)
// 			}

// 			if existing == nil {
// 				// Crear nuevo stock
// 				if err := s.repo.Create(ctx, stock); err != nil {
// 					s.logger.Error("Error creating stock",
// 						zap.String("ticker", item.Ticker),
// 						zap.Error(err))
// 					return count, fmt.Errorf("error creating stock: %w", err)
// 				}
// 			} else {
// 				// Truncar la fecha a segundos para evitar problemas de precisión
// 				timeValue = timeValue.Truncate(time.Second)
// 				existing.Time = existing.Time.Truncate(time.Second)

// 				// Actualizar stock existente si la fecha es más reciente
// 				if timeValue.After(existing.Time) {
// 					stock.ID = existing.ID
// 					if err := s.repo.Update(ctx, stock); err != nil {
// 						s.logger.Error("Error updating stock:",
// 							zap.String("ticker", item.Ticker),
// 							zap.Error(err))
// 						return count, fmt.Errorf("error updating stock: %w", err)
// 					}
// 				}
// 			}

// 			count++
// 		}

// 		// Verificar si hay más páginas
// 		nextPage = response.NextPage
// 		if nextPage == "" {
// 			break
// 		}

// 		// Añadir delay de 1 segundo entre peticiones
// 		//time.Sleep(1 * time.Second)
// 	}

// 	s.logger.Info("Synchronization completed", zap.Int("count", count), zap.Duration("duration", time.Since(timeStart)))
// 	return count, nil
// }

// GetRecommendations obtiene recomendaciones de stocks para invertir
func (s *stockService) GetRecommendations(ctx context.Context) ([]models.StockRecommendation, error) {
	// Obtener todos los stocks
	stocks, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Error getting stocks for recommendations", zap.Error(err))
		return nil, err
	}

	return s.processRecommendations(stocks), nil
}

// GetRecommendationsByTime obtiene recomendaciones de stocks por fecha
func (s *stockService) GetRecommendationsByTime(ctx context.Context, time string) ([]models.StockRecommendation, error) {
	// Obtener stocks por fecha
	stocks, err := s.repo.GetAllByTime(ctx, time)
	if err != nil {
		s.logger.Error("Error getting stocks for recommendations", zap.Error(err))
		return nil, err
	}

	return s.processRecommendations(stocks), nil
}

// processRecommendations procesa una lista de stocks y devuelve recomendaciones ordenadas
func (s *stockService) processRecommendations(stocks []models.Stock) []models.StockRecommendation {
	recommendations := make([]models.StockRecommendation, 0, len(stocks))

	for _, stock := range stocks {
		// Calcular puntuación y razón
		score, reason, potentialUp := s.calculateRecommendationScore(&stock)

		recommendations = append(recommendations, models.StockRecommendation{
			Stock:       stock,
			Score:       score,
			Reason:      reason,
			PotentialUp: potentialUp,
		})
	}

	// Ordenar recomendaciones por puntuación (de mayor a menor)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limitar a 5 recomendaciones
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations
}

// calculateRecommendationScore calcula una puntuación para una recomendación
func (s *stockService) calculateRecommendationScore(stock *models.Stock) (float64, []string, float64) {
	var score float64
	var reasons []string
	var potentialUp float64

	// Puntuación base por mejora de calificación
	if stock.RatingFrom != stock.RatingTo {
		ratingScore := s.getRatingScore(stock.RatingTo) - s.getRatingScore(stock.RatingFrom)
		score += ratingScore
		if ratingScore > 0 {
			reasons = append(reasons, fmt.Sprintf("Rating improvement from %s to %s", stock.RatingFrom, stock.RatingTo))
		}
	}

	// Puntuación por acción
	actionScore := s.getActionScore(stock.Action)
	score += actionScore
	if actionScore > 0 {
		reasons = append(reasons, fmt.Sprintf("Action taken %s %s", stock.Action, stock.Brokerage))
	}

	// Puntuación por potencial de crecimiento basado en target
	if stock.TargetFrom != "" && stock.TargetTo != "" {
		fromValue, fromErr := s.parsePrice(stock.TargetFrom)
		toValue, toErr := s.parsePrice(stock.TargetTo)

		if fromErr == nil && toErr == nil && fromValue > 0 {
			growthPercent := (toValue - fromValue) / fromValue * 100
			potentialUp = growthPercent

			// Puntuar si hay un aumento en el precio objetivo
			if growthPercent > 0 {
				score += growthPercent / 10 // Normalizar el impacto
				reasons = append(reasons, fmt.Sprintf("Increase in target price by %.2f%%", growthPercent))
			} else {
				reasons = append(reasons, fmt.Sprintf("Target price decreased by %.2f%%", growthPercent))
			}
		}
	}

	// Puntuación por broker
	// brokerScore := s.getBrokerScore(stock.Brokerage)
	// if brokerScore > 0 {
	// 	score += brokerScore
	// 	reasons = append(reasons, fmt.Sprintf("Recommendation of %s (highly trusted broker)", stock.Brokerage))
	// }

	// Razón final
	if len(reasons) == 0 {
		reasons = append(reasons, "Recommendation based on general analysis")
	}

	return score, reasons, potentialUp
}

// getRatingScore asigna una puntuación numérica a una calificación
func (s *stockService) getRatingScore(rating string) float64 {
	switch rating {
	case "Strong-Buy", "Buy", "Speculative Buy", "Positive", "Market Outperform", "Sector Outperform":
		return 5
	case "Outperform", "Outperformer", "Overweight":
		return 4
	case "Neutral", "Hold", "Equal Weight", "In-Line", "Inline", "Sector Perform", "Market Perform", "Sector Weight":
		return 3
	case "Underweight", "Underperform", "Sector Underperform", "Reduce", "Negative":
		return 2
	case "Sell":
		return 1
	default:
		return 0
	}
}

// getActionScore asigna una puntuación numérica a una acción
func (s *stockService) getActionScore(action string) float64 {
	switch action {
	case "upgraded by":
		return 8
	case "target raised by":
		return 7
	case "target set by":
		return 6
	case "initiated by":
		return 5
	case "reiterated by":
		return 4
	case "downgraded by":
		return 3
	case "target lowered by":
		return 1
	default:
		return 0
	}
}

// TODO: Implementar una puntuación real basada en un API como Alphavantage
// getBrokerScore asigna una puntuación de confianza a un broker
func (s *stockService) getBrokerScore(broker string) float64 {
	switch broker {
	case "The Goldman Sachs Group":
		return 3
	case "JP Morgan":
		return 2
	default:
		return 1
	}
}

// parsePrice convierte un string de precio a float64
func (s *stockService) parsePrice(price string) (float64, error) {
	// Eliminar el símbolo de dólar y espacios
	price = strings.TrimSpace(strings.Replace(price, "$", "", -1))
	return strconv.ParseFloat(price, 64)
}
