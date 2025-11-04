# ArbitraX

High-performance algorithmic trading simulator with backtesting, order book simulation, and strategy evaluation.

## Features

- **Order Book Engine** – Go-based price-time priority matching
- **Strategy Engine** – Python backtesting with pandas and backtrader
- **Live Simulation** – Real-time order matching and execution
- **Risk Analytics** – Position sizing, drawdown analysis, performance metrics
- **REST API** – Order management, strategy execution, market data
- **Dashboard** – React + TypeScript real-time visualization

## Tech Stack

- **Backend**: Go (Gin), Python (FastAPI)
- **Database**: PostgreSQL, Redis
- **Frontend**: React, TypeScript, Recharts
- **Infrastructure**: Docker, Terraform
- **CI/CD**: GitHub Actions

## Quick Start

```bash
# Start infrastructure
docker-compose up -d

# Verify services
docker-compose ps

# View logs
docker-compose logs -f
```

**Services:**
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- pgAdmin: `http://localhost:5050` (admin@arbitrax.local / admin)
- Redis Commander: `http://localhost:8081`

## Project Structure

```
arbitrax/
├── backend/          # Go order book + matching engine
├── strategy-engine/  # Python backtesting
├── frontend/         # React dashboard
├── infra/            # Terraform configs
└── docker-compose.yml
```


## License

MIT
