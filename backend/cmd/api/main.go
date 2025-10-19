package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/acagliol/arbitrax/backend/internal/matching"
	"github.com/acagliol/arbitrax/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

type OrderRequest struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Type     string  `json:"type" binding:"required,oneof=market limit stop_loss"`
	Side     string  `json:"side" binding:"required,oneof=buy sell"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Price    float64 `json:"price"` // Required for limit and stop_loss orders
}

type OrderResponse struct {
	Order  *models.Order   `json:"order"`
	Trades []*models.Trade `json:"trades,omitempty"`
}

var engine *matching.MatchingEngine

func main() {
	// Initialize matching engine
	engine = matching.NewMatchingEngine()

	// Create Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "arbitrax-backend",
		})
	})

	// Serve static frontend
	router.Static("/static", "../../frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("../../frontend/index.html")
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Order endpoints
		v1.POST("/orders", submitOrder)
		v1.GET("/orderbook/:symbol", getOrderBook)
		v1.GET("/trades/:symbol", getTrades)
	}

	// Start server
	router.Run(":8080")
}

// submitOrder handles order submission
func submitOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate price for limit and stop_loss orders
	if (req.Type == "limit" || req.Type == "stop_loss") && req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price is required for limit and stop_loss orders"})
		return
	}

	// Create order
	order := models.NewOrder(
		req.Symbol,
		models.OrderType(req.Type),
		models.OrderSide(req.Side),
		req.Quantity,
		req.Price,
	)

	// Submit to matching engine
	trades := engine.SubmitOrder(order)

	c.JSON(http.StatusOK, OrderResponse{
		Order:  order,
		Trades: trades,
	})
}

// getOrderBook returns the current order book for a symbol
func getOrderBook(c *gin.Context) {
	symbol := c.Param("symbol")

	ob := engine.GetOrderBook(symbol)
	if ob == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order book not found"})
		return
	}

	snapshot := ob.Snapshot()
	c.JSON(http.StatusOK, snapshot)
}

// getTrades returns recent trades for a symbol
func getTrades(c *gin.Context) {
	symbol := c.Param("symbol")

	// Get limit from query param (default 50, max 500)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 500 {
				limit = 500
			}
		}
	}

	trades := engine.GetRecentTrades(symbol, limit)
	c.JSON(http.StatusOK, gin.H{
		"symbol": symbol,
		"trades": trades,
		"count":  len(trades),
	})
}
