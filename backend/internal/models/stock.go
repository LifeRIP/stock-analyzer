package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Stock representa la información de una acción
type Stock struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid"`
	Ticker     string    `json:"ticker" gorm:"index;not null"`
	Company    string    `json:"company" gorm:"not null"`
	Brokerage  string    `json:"brokerage" gorm:"not null"`
	Action     string    `json:"action" gorm:"not null"`
	RatingFrom string    `json:"rating_from" gorm:"not null"`
	RatingTo   string    `json:"rating_to" gorm:"not null"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Time       time.Time `json:"time" gorm:"index;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Hook BeforeCreate se ejecuta antes de crear un registro
func (s *Stock) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

// StockResponse representa la respuesta de la API externa
type StockResponse struct {
	Items    []StockItem `json:"items"`
	NextPage string      `json:"next_page"`
}

// StockItem representa un elemento de stock en la respuesta de la API externa
type StockItem struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

// StockRecommendation representa una recomendación de inversión
type StockRecommendation struct {
	Stock       Stock    `json:"stock"`
	Score       float64  `json:"score"`
	Reason      []string `json:"reasons"`
	PotentialUp float64  `json:"potential_up,omitempty"`
}
