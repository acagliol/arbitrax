package orderbook

import (
	"container/heap"

	"github.com/acagliol/arbitrax/backend/internal/models"
)

// PriceLevel represents a price level in the order book with multiple orders
type PriceLevel struct {
	Price  float64
	Orders []*models.Order
}

// PriceLevelHeap is a heap of price levels
// For bids (buy orders), we want max-heap (highest price first)
// For asks (sell orders), we want min-heap (lowest price first)
type PriceLevelHeap struct {
	Levels []*PriceLevel
	IsBid  bool // true for bid (max-heap), false for ask (min-heap)
}

// Len returns the number of price levels
func (h *PriceLevelHeap) Len() int {
	return len(h.Levels)
}

// Less compares two price levels
func (h *PriceLevelHeap) Less(i, j int) bool {
	if h.IsBid {
		// For bids, higher price has priority (max-heap)
		return h.Levels[i].Price > h.Levels[j].Price
	}
	// For asks, lower price has priority (min-heap)
	return h.Levels[i].Price < h.Levels[j].Price
}

// Swap swaps two price levels
func (h *PriceLevelHeap) Swap(i, j int) {
	h.Levels[i], h.Levels[j] = h.Levels[j], h.Levels[i]
}

// Push adds a price level to the heap
func (h *PriceLevelHeap) Push(x interface{}) {
	h.Levels = append(h.Levels, x.(*PriceLevel))
}

// Pop removes and returns the top price level
func (h *PriceLevelHeap) Pop() interface{} {
	old := h.Levels
	n := len(old)
	level := old[n-1]
	h.Levels = old[0 : n-1]
	return level
}

// Peek returns the top price level without removing it
func (h *PriceLevelHeap) Peek() *PriceLevel {
	if h.Len() == 0 {
		return nil
	}
	return h.Levels[0]
}

// NewBidHeap creates a new max-heap for bid orders
func NewBidHeap() *PriceLevelHeap {
	h := &PriceLevelHeap{
		Levels: make([]*PriceLevel, 0),
		IsBid:  true,
	}
	heap.Init(h)
	return h
}

// NewAskHeap creates a new min-heap for ask orders
func NewAskHeap() *PriceLevelHeap {
	h := &PriceLevelHeap{
		Levels: make([]*PriceLevel, 0),
		IsBid:  false,
	}
	heap.Init(h)
	return h
}

// AddOrder adds an order to the appropriate price level
func (h *PriceLevelHeap) AddOrder(order *models.Order) {
	// Find existing price level
	for _, level := range h.Levels {
		if level.Price == order.Price {
			level.Orders = append(level.Orders, order)
			return
		}
	}

	// Create new price level
	newLevel := &PriceLevel{
		Price:  order.Price,
		Orders: []*models.Order{order},
	}
	heap.Push(h, newLevel)
}

// RemoveOrder removes an order from the heap
func (h *PriceLevelHeap) RemoveOrder(order *models.Order) bool {
	for i, level := range h.Levels {
		if level.Price == order.Price {
			for j, o := range level.Orders {
				if o.ID == order.ID {
					// Remove order from price level
					level.Orders = append(level.Orders[:j], level.Orders[j+1:]...)

					// If price level is empty, remove it
					if len(level.Orders) == 0 {
						h.Levels = append(h.Levels[:i], h.Levels[i+1:]...)
						heap.Init(h) // Re-heapify
					}
					return true
				}
			}
		}
	}
	return false
}
