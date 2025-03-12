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
		return nil, fmt.Errorf("error creating request to external API: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	maxAttempts := 5
	var lastErr error

	for i := 1; i <= maxAttempts; i++ {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error making request to external API (attempt %d/%d): %w",
				i, maxAttempts, err)

			return nil, lastErr
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("error external API responded with a status code%d (attempt %d/%d)",
				resp.StatusCode, i, maxAttempts)

			// Reintentar si la API externa responde con un error
			if i < maxAttempts {
				c.logger.Warn("Retrying in 5 seconds...",
					zap.Int("attempt", i),
					zap.Int("max_attempts", maxAttempts),
					zap.Error(lastErr))
				time.Sleep(5 * time.Second)
				continue
			}
			return nil, lastErr
		}

		var stockResponse models.StockResponse
		if err := json.NewDecoder(resp.Body).Decode(&stockResponse); err != nil {
			lastErr = fmt.Errorf("error decoding API response (attempt %d/%d): %w",
				i, maxAttempts, err)
			return nil, lastErr
		}

		return &stockResponse, nil
	}
	return nil, fmt.Errorf("Unexpected error after %d attempts", maxAttempts)
}
