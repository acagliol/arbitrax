package matching

import (
	"container/heap"
	"sync"

	"github.com/acagliol/arbitrax/backend/internal/models"
	"github.com/acagliol/arbitrax/backend/internal/orderbook"
)

// MatchingEngine handles order matching across multiple order books
type MatchingEngine struct {
	orderBooks map[string]*orderbook.OrderBook
	trades     []*models.Trade
	mutex      sync.RWMutex
}

// NewMatchingEngine creates a new matching engine
func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		orderBooks: make(map[string]*orderbook.OrderBook),
		trades:     make([]*models.Trade, 0),
	}
}

// GetOrCreateOrderBook gets or creates an order book for a symbol
func (me *MatchingEngine) GetOrCreateOrderBook(symbol string) *orderbook.OrderBook {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if ob, exists := me.orderBooks[symbol]; exists {
		return ob
	}

	ob := orderbook.NewOrderBook(symbol)
	me.orderBooks[symbol] = ob
	return ob
}

// GetOrderBook retrieves an order book for a symbol
func (me *MatchingEngine) GetOrderBook(symbol string) *orderbook.OrderBook {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me.orderBooks[symbol]
}

// SubmitOrder submits an order to the matching engine
func (me *MatchingEngine) SubmitOrder(order *models.Order) []*models.Trade {
	ob := me.GetOrCreateOrderBook(order.Symbol)

	var trades []*models.Trade

	// Handle different order types
	switch order.Type {
	case models.OrderTypeMarket:
		trades = me.matchMarketOrder(ob, order)
	case models.OrderTypeLimit:
		trades = me.matchLimitOrder(ob, order)
	case models.OrderTypeStopLoss:
		// Stop-loss orders become market orders when triggered
		// For now, we'll treat them as limit orders at the stop price
		order.Type = models.OrderTypeLimit
		trades = me.matchLimitOrder(ob, order)
	}

	// Store trades
	if len(trades) > 0 {
		me.mutex.Lock()
		me.trades = append(me.trades, trades...)
		me.mutex.Unlock()
	}

	return trades
}

// matchMarketOrder matches a market order immediately at best available prices
func (me *MatchingEngine) matchMarketOrder(ob *orderbook.OrderBook, order *models.Order) []*models.Trade {
	trades := make([]*models.Trade, 0)

	var oppositeHeap *orderbook.PriceLevelHeap
	if order.Side == models.OrderSideBuy {
		oppositeHeap = ob.Asks
	} else {
		oppositeHeap = ob.Bids
	}

	// Match against all available opposite orders until filled
	for order.RemainingQuantity() > 0 && oppositeHeap.Len() > 0 {
		bestLevel := oppositeHeap.Peek()
		if bestLevel == nil {
			break
		}
		if len(bestLevel.Orders) == 0 {
			heap.Pop(oppositeHeap)
			continue
		}

		// Match with orders at this price level (FIFO - time priority)
		for len(bestLevel.Orders) > 0 && order.RemainingQuantity() > 0 {
			oppositeOrder := bestLevel.Orders[0]

			// Calculate trade quantity
			tradeQty := min(order.RemainingQuantity(), oppositeOrder.RemainingQuantity())
			tradePrice := oppositeOrder.Price

			// Create trade
			var trade *models.Trade
			if order.Side == models.OrderSideBuy {
				trade = models.NewTrade(order.Symbol, order.ID, oppositeOrder.ID, tradePrice, tradeQty)
			} else {
				trade = models.NewTrade(order.Symbol, oppositeOrder.ID, order.ID, tradePrice, tradeQty)
			}

			// Fill both orders
			order.Fill(tradeQty, tradePrice)
			oppositeOrder.Fill(tradeQty, tradePrice)

			// Update last price
			ob.LastPrice = tradePrice
			ob.LastTrade = trade

			trades = append(trades, trade)

			// If opposite order is filled, remove it from the book
			if oppositeOrder.IsFilled() {
				bestLevel.Orders = bestLevel.Orders[1:]
			}

			// If incoming order is filled, stop matching at this level
			if order.IsFilled() {
				break
			}
		}

		// If price level is empty, remove it
		if len(bestLevel.Orders) == 0 {
			heap.Pop(oppositeHeap)
		}
	}

	return trades
}

// matchLimitOrder matches a limit order, adding remainder to order book if not fully filled
func (me *MatchingEngine) matchLimitOrder(ob *orderbook.OrderBook, order *models.Order) []*models.Trade {
	trades := make([]*models.Trade, 0)

	var oppositeHeap *orderbook.PriceLevelHeap
	if order.Side == models.OrderSideBuy {
		oppositeHeap = ob.Asks
	} else {
		oppositeHeap = ob.Bids
	}

	// Match against opposite orders while price is acceptable
	for order.RemainingQuantity() > 0 && oppositeHeap.Len() > 0 {
		bestLevel := oppositeHeap.Peek()
		if bestLevel == nil || len(bestLevel.Orders) == 0 {
			break
		}

		// Check if price is acceptable
		if order.Side == models.OrderSideBuy && bestLevel.Price > order.Price {
			break // Ask price too high
		}
		if order.Side == models.OrderSideSell && bestLevel.Price < order.Price {
			break // Bid price too low
		}

		// Match with orders at this price level (FIFO - time priority)
		for len(bestLevel.Orders) > 0 && order.RemainingQuantity() > 0 {
			oppositeOrder := bestLevel.Orders[0]

			// Calculate trade quantity
			tradeQty := min(order.RemainingQuantity(), oppositeOrder.RemainingQuantity())
			tradePrice := oppositeOrder.Price

			// Create trade
			var trade *models.Trade
			if order.Side == models.OrderSideBuy {
				trade = models.NewTrade(order.Symbol, order.ID, oppositeOrder.ID, tradePrice, tradeQty)
			} else {
				trade = models.NewTrade(order.Symbol, oppositeOrder.ID, order.ID, tradePrice, tradeQty)
			}

			// Fill both orders
			order.Fill(tradeQty, tradePrice)
			oppositeOrder.Fill(tradeQty, tradePrice)

			// Update last price
			ob.LastPrice = tradePrice
			ob.LastTrade = trade

			trades = append(trades, trade)

			// If opposite order is filled, remove it
			if oppositeOrder.IsFilled() {
				bestLevel.Orders = bestLevel.Orders[1:]
			}
		}

		// If price level is empty, remove it
		if len(bestLevel.Orders) == 0 {
			heap.Pop(oppositeHeap)
		}
	}

	// If order is not fully filled, add remainder to order book
	if order.RemainingQuantity() > 0 {
		ob.AddOrder(order)
	}

	return trades
}

// GetRecentTrades returns recent trades for a symbol
func (me *MatchingEngine) GetRecentTrades(symbol string, limit int) []*models.Trade {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	result := make([]*models.Trade, 0)
	count := 0

	// Iterate from most recent
	for i := len(me.trades) - 1; i >= 0 && count < limit; i-- {
		if me.trades[i].Symbol == symbol {
			result = append(result, me.trades[i])
			count++
		}
	}

	return result
}

// Helper function to get minimum of two floats
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
