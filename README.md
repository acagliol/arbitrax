# ArbitraX

> High-performance algorithmic trading simulator with backtesting, order book simulation, and real-time paper trading

## Overview

ArbitraX is a comprehensive trading sandbox that combines:
- **Order Book Engine** (Go) - Low-latency order matching with price-time priority
- **Backtesting Framework** (Python) - Strategy evaluation on historical data
- **Risk Analytics** - Sharpe ratio, VaR, drawdown analysis, Monte Carlo simulation
- **Real-time Paper Trading** - Live strategy execution with simulated orders
- **Interactive Dashboard** - React frontend with charts and performance metrics

## Tech Stack

- **Backend:** Go 1.21+ (Gin framework)
- **Strategy Engine:** Python 3.11+ (pandas, backtrader, scikit-learn)
- **Database:** PostgreSQL 15 + Redis 7
- **Frontend:** React 18 + TypeScript + Recharts
- **Infrastructure:** Docker, Terraform, AWS

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Python 3.11+
- Node.js 18+

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/arbitrax.git
cd arbitrax
```

2. Start infrastructure:
```bash
docker-compose up -d
```

3. Run database migrations:
```bash
cd infra/database/migrations
./migrate.sh up
```

4. Start the Go backend:
```bash
cd backend
go run cmd/api/main.go
```

5. Start the Python strategy engine:
```bash
cd strategy-engine
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py
```

6. Start the frontend:
```bash
cd frontend
npm install
npm run dev
```

Visit http://localhost:3000 to access the dashboard.

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Frontend  │────▶│  API Gateway │────▶│   Backend   │
│   (React)   │     │   (Go/Gin)   │     │ (Order Book)│
└─────────────┘     └──────────────┘     └─────────────┘
                            │                     │
                            ▼                     ▼
                    ┌──────────────┐     ┌─────────────┐
                    │   Strategy   │     │   Redis     │
                    │ Engine (Py)  │     │(Live State) │
                    └──────────────┘     └─────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │  PostgreSQL  │
                    │(Historical)  │
                    └──────────────┘
```

## Project Structure

```
arbitrax/
├── backend/              # Go order book & matching engine
│   ├── cmd/api/         # Main application entry
│   ├── internal/        # Private application code
│   │   ├── orderbook/   # Order book implementation
│   │   ├── matching/    # Order matching engine
│   │   └── models/      # Data models
│   └── pkg/             # Public libraries
├── strategy-engine/      # Python backtesting framework
│   ├── backtesting/     # Backtesting engine
│   ├── strategies/      # Trading strategies
│   ├── data/            # Data ingestion & processing
│   └── utils/           # Helper functions
├── api/                  # API gateway (if separate from backend)
├── frontend/             # React dashboard
│   └── src/
│       ├── components/  # UI components
│       ├── pages/       # Page components
│       └── api/         # API client
├── infra/                # Infrastructure as code
│   ├── terraform/       # Terraform configs
│   ├── docker/          # Dockerfiles
│   └── database/        # Database schemas & migrations
└── docker-compose.yml    # Local development setup
```

## API Endpoints

### Order Management
- `POST /api/v1/orders` - Submit new order
- `GET /api/v1/orders/:id` - Get order details
- `GET /api/v1/orderbook/:symbol` - Get current order book
- `GET /api/v1/trades/:symbol` - Get recent trades

### Backtesting
- `POST /api/v1/backtests` - Start new backtest
- `GET /api/v1/backtests/:id` - Get backtest results
- `GET /api/v1/backtests/:id/trades` - Get trade log

### Strategies
- `GET /api/v1/strategies` - List available strategies
- `POST /api/v1/strategies` - Create custom strategy
- `GET /api/v1/strategies/:id` - Get strategy details

### Assets
- `GET /api/v1/assets` - List available assets
- `GET /api/v1/assets/:symbol/prices` - Get historical prices

## Development

### Running Tests

Go backend:
```bash
cd backend
go test ./... -v -cover
```

Python strategy engine:
```bash
cd strategy-engine
pytest tests/ -v --cov
```

### Database Migrations

Create new migration:
```bash
migrate create -ext sql -dir infra/database/migrations -seq create_users_table
```

Run migrations:
```bash
migrate -path infra/database/migrations -database "postgresql://user:pass@localhost:5432/arbitrax?sslmode=disable" up
```

## Roadmap

- [x] Week 1: Architecture & infrastructure setup
- [ ] Week 2: Order book & matching engine
- [ ] Week 3: Historical data & backtesting
- [ ] Week 4: API integration & basic frontend
- [ ] Week 5: Risk analytics & strategy comparison
- [ ] Week 6: Live mode & real-time features
- [ ] Week 7: ML strategies & optimization
- [ ] Week 8: Production deployment

See [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md) for detailed timeline.

## Contributing

This is a portfolio project, but suggestions and feedback are welcome! Feel free to open issues or submit PRs.

## License

MIT License - see [LICENSE](./LICENSE) for details.

## Author

Built by [Your Name] as a portfolio project demonstrating systems engineering, financial modeling, and full-stack development.

## Acknowledgments

- Inspired by QuantConnect and Interactive Brokers
- Data provided by Alpha Vantage and Yahoo Finance
- Built with support from the open-source community
