package orderbook

import (
	"sync"
	"time"

	"github.com/acagliol/arbitrax/backend/internal/models"
	"github.com/google/uuid"
)

// OrderBook represents the order book for a single symbol
type OrderBook struct {
	Symbol    string
	Bids      *PriceLevelHeap
	Asks      *PriceLevelHeap
	LastPrice float64
	LastTrade *models.Trade
	Timestamp time.Time
	mutex     sync.RWMutex
	orders    map[uuid.UUID]*models.Order // Track all orders by ID
}

// NewOrderBook creates a new order book for a symbol
func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{
		Symbol:    symbol,
		Bids:      NewBidHeap(),
		Asks:      NewAskHeap(),
		LastPrice: 0,
		Timestamp: time.Now(),
		orders:    make(map[uuid.UUID]*models.Order),
	}
}

// AddOrder adds an order to the order book
func (ob *OrderBook) AddOrder(order *models.Order) {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	// Store order
	ob.orders[order.ID] = order

	// Add to appropriate side
	if order.Side == models.OrderSideBuy {
		ob.Bids.AddOrder(order)
	} else {
		ob.Asks.AddOrder(order)
	}

	ob.Timestamp = time.Now()
}

// RemoveOrder removes an order from the order book
func (ob *OrderBook) RemoveOrder(orderID uuid.UUID) bool {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	order, exists := ob.orders[orderID]
	if !exists {
		return false
	}

	delete(ob.orders, orderID)

	if order.Side == models.OrderSideBuy {
		return ob.Bids.RemoveOrder(order)
	}
	return ob.Asks.RemoveOrder(order)
}

// GetOrder retrieves an order by ID
func (ob *OrderBook) GetOrder(orderID uuid.UUID) (*models.Order, bool) {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	order, exists := ob.orders[orderID]
	return order, exists
}

// GetBestBid returns the highest bid price
func (ob *OrderBook) GetBestBid() float64 {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	if ob.Bids.Len() == 0 {
		return 0
	}
	return ob.Bids.Peek().Price
}

// GetBestAsk returns the lowest ask price
func (ob *OrderBook) GetBestAsk() float64 {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	if ob.Asks.Len() == 0 {
		return 0
	}
	return ob.Asks.Peek().Price
}

// GetSpread returns the bid-ask spread
func (ob *OrderBook) GetSpread() float64 {
	bestBid := ob.GetBestBid()
	bestAsk := ob.GetBestAsk()

	if bestBid == 0 || bestAsk == 0 {
		return 0
	}

	return bestAsk - bestBid
}

// GetMidPrice returns the mid-market price
func (ob *OrderBook) GetMidPrice() float64 {
	bestBid := ob.GetBestBid()
	bestAsk := ob.GetBestAsk()

	if bestBid == 0 || bestAsk == 0 {
		return ob.LastPrice
	}

	return (bestBid + bestAsk) / 2
}

// Snapshot returns a snapshot of the order book
func (ob *OrderBook) Snapshot() *OrderBookSnapshot {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	snapshot := &OrderBookSnapshot{
		Symbol:    ob.Symbol,
		Bids:      make([]PriceLevelSnapshot, 0),
		Asks:      make([]PriceLevelSnapshot, 0),
		LastPrice: ob.LastPrice,
		Timestamp: ob.Timestamp,
	}

	// Copy bid levels
	for _, level := range ob.Bids.Levels {
		totalQty := 0.0
		for _, order := range level.Orders {
			totalQty += order.RemainingQuantity()
		}
		snapshot.Bids = append(snapshot.Bids, PriceLevelSnapshot{
			Price:    level.Price,
			Quantity: totalQty,
			Orders:   len(level.Orders),
		})
	}

	// Copy ask levels
	for _, level := range ob.Asks.Levels {
		totalQty := 0.0
		for _, order := range level.Orders {
			totalQty += order.RemainingQuantity()
		}
		snapshot.Asks = append(snapshot.Asks, PriceLevelSnapshot{
			Price:    level.Price,
			Quantity: totalQty,
			Orders:   len(level.Orders),
		})
	}

	return snapshot
}

// OrderBookSnapshot is a read-only snapshot of the order book
type OrderBookSnapshot struct {
	Symbol    string                `json:"symbol"`
	Bids      []PriceLevelSnapshot  `json:"bids"`
	Asks      []PriceLevelSnapshot  `json:"asks"`
	LastPrice float64               `json:"last_price"`
	Timestamp time.Time             `json:"timestamp"`
}

// PriceLevelSnapshot represents a price level in the snapshot
type PriceLevelSnapshot struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Orders   int     `json:"orders"`
}
