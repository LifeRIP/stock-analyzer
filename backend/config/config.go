package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// Module proporciona las dependencias de la configuración
var Module = fx.Provide(LoadConfig)

// Config contiene la configuración de la aplicación
type Config struct {
	DatabaseURL string
	APIEndpoint string
	APIKey      string
	ServerPort  string
	Environment string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() (*Config, error) {
	// Cargar variables de entorno desde .env si existe
	_ = godotenv.Load()

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://root@localhost:26257/defaultdb?sslmode=disable"),
		APIEndpoint: getEnv("API_ENDPOINT", "http://localhost:8081"),
		APIKey:      getEnv("API_KEY", "exampleApiKey"),
		ServerPort:  getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}, nil
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
