package orderbook

import (
	"testing"

	"github.com/acagliol/arbitrax/backend/internal/models"
)

func TestNewOrderBook(t *testing.T) {
	ob := NewOrderBook("AAPL")

	if ob.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", ob.Symbol)
	}

	if ob.Bids == nil || ob.Asks == nil {
		t.Error("Bids or Asks heap is nil")
	}

	if ob.LastPrice != 0 {
		t.Errorf("Expected last price 0, got %f", ob.LastPrice)
	}
}

func TestAddOrder(t *testing.T) {
	ob := NewOrderBook("AAPL")

	buyOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0)
	sellOrder := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 151.0)

	ob.AddOrder(buyOrder)
	ob.AddOrder(sellOrder)

	if ob.Bids.Len() != 1 {
		t.Errorf("Expected 1 bid level, got %d", ob.Bids.Len())
	}

	if ob.Asks.Len() != 1 {
		t.Errorf("Expected 1 ask level, got %d", ob.Asks.Len())
	}
}

func TestGetBestBidAsk(t *testing.T) {
	ob := NewOrderBook("AAPL")

	// Add multiple bid orders
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 149.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 151.0))

	// Add multiple ask orders
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 153.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 152.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 154.0))

	bestBid := ob.GetBestBid()
	if bestBid != 151.0 {
		t.Errorf("Expected best bid 151.0, got %f", bestBid)
	}

	bestAsk := ob.GetBestAsk()
	if bestAsk != 152.0 {
		t.Errorf("Expected best ask 152.0, got %f", bestAsk)
	}
}

func TestGetSpread(t *testing.T) {
	ob := NewOrderBook("AAPL")

	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 152.0))

	spread := ob.GetSpread()
	expected := 2.0
	if spread != expected {
		t.Errorf("Expected spread %f, got %f", expected, spread)
	}
}

func TestGetMidPrice(t *testing.T) {
	ob := NewOrderBook("AAPL")

	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 100, 152.0))

	midPrice := ob.GetMidPrice()
	expected := 151.0
	if midPrice != expected {
		t.Errorf("Expected mid price %f, got %f", expected, midPrice)
	}
}

func TestRemoveOrder(t *testing.T) {
	ob := NewOrderBook("AAPL")

	order := models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0)
	ob.AddOrder(order)

	if ob.Bids.Len() != 1 {
		t.Errorf("Expected 1 bid level, got %d", ob.Bids.Len())
	}

	removed := ob.RemoveOrder(order.ID)
	if !removed {
		t.Error("Failed to remove order")
	}

	if ob.Bids.Len() != 0 {
		t.Errorf("Expected 0 bid levels, got %d", ob.Bids.Len())
	}
}

func TestSnapshot(t *testing.T) {
	ob := NewOrderBook("AAPL")

	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 100, 150.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideBuy, 50, 150.0))
	ob.AddOrder(models.NewOrder("AAPL", models.OrderTypeLimit, models.OrderSideSell, 75, 152.0))

	snapshot := ob.Snapshot()

	if snapshot.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", snapshot.Symbol)
	}

	if len(snapshot.Bids) != 1 {
		t.Errorf("Expected 1 bid level in snapshot, got %d", len(snapshot.Bids))
	}

	if snapshot.Bids[0].Quantity != 150.0 {
		t.Errorf("Expected bid quantity 150.0, got %f", snapshot.Bids[0].Quantity)
	}

	if snapshot.Bids[0].Orders != 2 {
		t.Errorf("Expected 2 orders at bid level, got %d", snapshot.Bids[0].Orders)
	}
}
