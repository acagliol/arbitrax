package matching

import (
	"testing"

	"github.com/acagliol/arbitrax/backend/internal/models"
)

func TestNewMatchingEngine(t *testing.T) {
	me := NewMatchingEngine()

	if me == nil {
		t.Fatal("NewMatchingEngine returned nil")
	}

	if me.orderBooks == nil {
		t.Error("orderBooks map is nil")
	}
}

func TestGetOrCreateOrderBook(t *testing.T) {
	me := NewMatchingEngine()

	ob1 := me.GetOrCreateOrderBook("AAPL")
	if ob1 == nil {
		t.Fatal("GetOrCreateOrderBook returned nil")
	}

	if ob1.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", ob1.Symbol)
	}

	// Get the same order book again
	ob2 := me.GetOrCreateOrderBook("AAPL")
	if ob1 != ob2 {
		t.Error("GetOrCreateOrderBook should return the same instance")
	}
}

func TestMatchLimitOrders(t *testing.T) {
	me := NewMatchingEngine()

	// Add sell order
	sellOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 150.0)
	trades := me.SubmitOrder(sellOrder)

	if len(trades) != 0 {
		t.Errorf("Expected no trades for first order, got %d", len(trades))
	}

	// Add matching buy order
	buyOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0)
	trades = me.SubmitOrder(buyOrder)

	if len(trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(trades))
	}

	trade := trades[0]
	if trade.Price != 150.0 {
		t.Errorf("Expected trade price 150.0, got %f", trade.Price)
	}

	if trade.Quantity != 100 {
		t.Errorf("Expected trade quantity 100, got %f", trade.Quantity)
	}

	// Check that both orders are filled
	if !buyOrder.IsFilled() {
		t.Error("Buy order should be filled")
	}

	if !sellOrder.IsFilled() {
		t.Error("Sell order should be filled")
	}
}

func TestPartialFill(t *testing.T) {
	me := NewMatchingEngine()

	// Add sell order for 100 shares
	sellOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 150.0)
	me.SubmitOrder(sellOrder)

	// Add buy order for 50 shares (partial fill)
	buyOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 50, 150.0)
	trades := me.SubmitOrder(buyOrder)

	if len(trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(trades))
	}

	if trades[0].Quantity != 50 {
		t.Errorf("Expected trade quantity 50, got %f", trades[0].Quantity)
	}

	if !buyOrder.IsFilled() {
		t.Error("Buy order should be fully filled")
	}

	if sellOrder.IsFilled() {
		t.Error("Sell order should be partially filled")
	}

	if sellOrder.RemainingQuantity() != 50 {
		t.Errorf("Expected remaining quantity 50, got %f", sellOrder.RemainingQuantity())
	}
}

func TestMarketOrder(t *testing.T) {
	me := NewMatchingEngine()

	// Add sell limit orders at different prices
	me.SubmitOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 50, 150.0))
	me.SubmitOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 50, 151.0))

	// Add market buy order for 100 shares
	marketOrder := models.NewOrder("AAPL", models.OrderTypeMarket, models.OrderSideBuy, 100, 0)
	trades := me.SubmitOrder(marketOrder)

	if len(trades) != 2 {
		t.Fatalf("Expected 2 trades, got %d", len(trades))
	}

	// First trade should be at 150.0 (best price)
	if trades[0].Price != 150.0 {
		t.Errorf("Expected first trade at 150.0, got %f", trades[0].Price)
	}

	// Second trade should be at 151.0
	if trades[1].Price != 151.0 {
		t.Errorf("Expected second trade at 151.0, got %f", trades[1].Price)
	}

	if !marketOrder.IsFilled() {
		t.Error("Market order should be fully filled")
	}
}

func TestPriceTimePriority(t *testing.T) {
	me := NewMatchingEngine()

	// Add two sell orders at the same price
	sellOrder1 := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 50, 150.0)
	sellOrder2 := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 50, 150.0)

	me.SubmitOrder(sellOrder1)
	me.SubmitOrder(sellOrder2)

	// Add buy order - should match with first sell order (time priority)
	buyOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 50, 150.0)
	trades := me.SubmitOrder(buyOrder)

	if len(trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(trades))
	}

	// First order should be filled
	if !sellOrder1.IsFilled() {
		t.Error("First sell order should be filled (time priority)")
	}

	// Second order should not be filled
	if sellOrder2.IsFilled() {
		t.Error("Second sell order should not be filled")
	}
}

func TestNoCrossing(t *testing.T) {
	me := NewMatchingEngine()

	// Add sell order at 152.0
	sellOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 152.0)
	me.SubmitOrder(sellOrder)

	// Add buy order at 150.0 (below sell price, should not match)
	buyOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0)
	trades := me.SubmitOrder(buyOrder)

	if len(trades) != 0 {
		t.Errorf("Expected no trades (no price crossing), got %d", len(trades))
	}

	// Both orders should be in the order book
	ob := me.GetOrderBook("AAPL")
	if ob.Bids.Len() != 1 || ob.Asks.Len() != 1 {
		t.Error("Both orders should be in the order book")
	}
}

func TestGetRecentTrades(t *testing.T) {
	me := NewMatchingEngine()

	// Execute some trades
	me.SubmitOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 150.0))
	me.SubmitOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0))

	me.SubmitOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 50, 151.0))
	me.SubmitOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 50, 151.0))

	// Get recent trades
	trades := me.GetRecentTrades("AAPL", 10)

	if len(trades) != 2 {
		t.Errorf("Expected 2 trades, got %d", len(trades))
	}

	// Most recent trade should be first
	if trades[0].Price != 151.0 {
		t.Errorf("Expected most recent trade at 151.0, got %f", trades[0].Price)
	}
}

func TestEmptyOrderBook(t *testing.T) {
	me := NewMatchingEngine()

	// Try to get order book that doesn't exist
	ob := me.GetOrderBook("NONEXISTENT")
	if ob != nil {
		t.Error("Expected nil for non-existent order book")
	}
}
