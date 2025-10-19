package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeMarket   OrderType = "market"
	OrderTypeLimit    OrderType = "limit"
	OrderTypeStopLoss OrderType = "stop_loss"
)

// OrderSide represents buy or sell
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPartial   OrderStatus = "partial"
	OrderStatusFilled    OrderStatus = "filled"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a trading order
type Order struct {
	ID             uuid.UUID   `json:"id"`
	Symbol         string      `json:"symbol"`
	Type           OrderType   `json:"type"`
	Side           OrderSide   `json:"side"`
	Quantity       float64     `json:"quantity"`
	Price          float64     `json:"price"` // 0 for market orders
	Status         OrderStatus `json:"status"`
	FilledQuantity float64     `json:"filled_quantity"`
	FilledPrice    float64     `json:"filled_price"`
	SubmittedAt    time.Time   `json:"submitted_at"`
	FilledAt       *time.Time  `json:"filled_at,omitempty"`
	CancelledAt    *time.Time  `json:"cancelled_at,omitempty"`
}

// NewOrder creates a new order
func NewOrder(symbol string, orderType OrderType, side OrderSide, quantity, price float64) *Order {
	return &Order{
		ID:             uuid.New(),
		Symbol:         symbol,
		Type:           orderType,
		Side:           side,
		Quantity:       quantity,
		Price:          price,
		Status:         OrderStatusPending,
		FilledQuantity: 0,
		FilledPrice:    0,
		SubmittedAt:    time.Now(),
	}
}

// RemainingQuantity returns the unfilled quantity
func (o *Order) RemainingQuantity() float64 {
	return o.Quantity - o.FilledQuantity
}

// IsFilled returns true if the order is completely filled
func (o *Order) IsFilled() bool {
	return o.FilledQuantity >= o.Quantity
}

// Fill partially or fully fills the order
func (o *Order) Fill(quantity, price float64) {
	o.FilledQuantity += quantity
	// Update filled price as weighted average
	if o.FilledQuantity > 0 {
		o.FilledPrice = ((o.FilledPrice * (o.FilledQuantity - quantity)) + (price * quantity)) / o.FilledQuantity
	}

	if o.IsFilled() {
		o.Status = OrderStatusFilled
		now := time.Now()
		o.FilledAt = &now
	} else if o.FilledQuantity > 0 {
		o.Status = OrderStatusPartial
	}
}
