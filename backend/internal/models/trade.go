package models

import (
	"time"

	"github.com/google/uuid"
)

// Trade represents an executed trade between a buy and sell order
type Trade struct {
	ID          uuid.UUID `json:"id"`
	Symbol      string    `json:"symbol"`
	BuyOrderID  uuid.UUID `json:"buy_order_id"`
	SellOrderID uuid.UUID `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
}

// NewTrade creates a new trade
func NewTrade(symbol string, buyOrderID, sellOrderID uuid.UUID, price, quantity float64) *Trade {
	return &Trade{
		ID:          uuid.New(),
		Symbol:      symbol,
		BuyOrderID:  buyOrderID,
		SellOrderID: sellOrderID,
		Price:       price,
		Quantity:    quantity,
		Timestamp:   time.Now(),
	}
}
