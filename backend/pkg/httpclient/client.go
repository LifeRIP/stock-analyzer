package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/liferip/stock-analyzer/backend/config"
	"github.com/liferip/stock-analyzer/backend/internal/models"
)

// Module proporciona las dependencias del cliente HTTP
var Module = fx.Provide(NewStockClient)

// StockClient es un cliente para la API de stocks
type StockClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	logger     *zap.Logger
}

// NewStockClient crea un nuevo cliente para la API de stocks
func NewStockClient(cfg *config.Config, logger *zap.Logger) *StockClient {
	return &StockClient{
		BaseURL: cfg.APIEndpoint,
		APIKey:  cfg.APIKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger.Named("stock_client"),
	}
}

// GetStocks obtiene la lista de stocks desde la API externa
func (c *StockClient) GetStocks(nextPage string) (*models.StockResponse, error) {
	url := c.BaseURL
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la solicitud a la API externa: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	maxIntentos := 3
	var lastErr error

	for i := 1; i <= maxIntentos; i++ {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error al realizar la solicitud a la API externa (intento %d/%d): %w",
				i, maxIntentos, err)

			return nil, lastErr
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("la API respondió con código de estado %d (intento %d/%d)",
				resp.StatusCode, i, maxIntentos)

			// Reintentar si la API externa responde con un error
			if i < maxIntentos {
				c.logger.Warn("Reintentando en 5 segundos...",
					zap.Int("intento", i),
					zap.Int("max_intentos", maxIntentos),
					zap.Error(lastErr))
				time.Sleep(5 * time.Second)
				continue
			}
			return nil, lastErr
		}

		var stockResponse models.StockResponse
		if err := json.NewDecoder(resp.Body).Decode(&stockResponse); err != nil {
			lastErr = fmt.Errorf("error al decodificar la respuesta de la API (intento %d/%d): %w",
				i, maxIntentos, err)
			return nil, lastErr
		}

		return &stockResponse, nil
	}
	return nil, fmt.Errorf("error inesperado después de %d intentos", maxIntentos)
}
