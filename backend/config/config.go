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
	DatabaseUser    string
	DatabasePass    string
	DatabaseHost    string
	DatabasePort    string
	DatabaseName    string
	DatabaseSSLMode string
	DatabaseURL     string
	APIEndpoint     string
	APIKey          string
	ServerPort      string
	Environment     string
	SwaggerHost     string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() (*Config, error) {
	// Cargar variables de entorno desde .env si existe
	_ = godotenv.Load()

	return &Config{
		DatabaseUser:    getEnv("DATABASE_USER", "root"),
		DatabasePass:    getEnv("DATABASE_PASS", ""),
		DatabaseHost:    getEnv("DATABASE_HOST", "cockroach"),
		DatabasePort:    getEnv("DATABASE_PORT", "26257"),
		DatabaseName:    getEnv("DATABASE_NAME", "stock_analyzer_db"),
		DatabaseSSLMode: getEnv("DATABASE_SSL_MODE", "disable"),
		APIEndpoint:     getEnv("API_ENDPOINT", "http://localhost:8083"),
		APIKey:          getEnv("API_KEY", "exampleApiKey"),
		ServerPort:      getEnv("PORT", "8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		SwaggerHost:     getEnv("SWAGGER_HOST", "localhost:8080"),
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
