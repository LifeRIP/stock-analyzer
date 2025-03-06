package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"

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
}

// NewStockClient crea un nuevo cliente para la API de stocks
func NewStockClient(cfg *config.Config) *StockClient {
	return &StockClient{
		BaseURL: cfg.APIEndpoint,
		APIKey:  cfg.APIKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la solicitud a la API externa: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("la API respondió con código de estado %d", resp.StatusCode)
	}

	var stockResponse models.StockResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockResponse); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta de la API: %w", err)
	}

	return &stockResponse, nil
}
