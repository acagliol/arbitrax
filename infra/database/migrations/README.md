# Database Migrations

This directory contains SQL migrations for the ArbitraX database.

## Migration Tool

We use [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

### Installation

```bash
# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate

# Or using Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Usage

Apply all migrations:
```bash
migrate -path ./infra/database/migrations \
  -database "postgresql://arbitrax:password@localhost:5432/arbitrax?sslmode=disable" \
  up
```

Rollback last migration:
```bash
migrate -path ./infra/database/migrations \
  -database "postgresql://arbitrax:password@localhost:5432/arbitrax?sslmode=disable" \
  down 1
```

Check migration version:
```bash
migrate -path ./infra/database/migrations \
  -database "postgresql://arbitrax:password@localhost:5432/arbitrax?sslmode=disable" \
  version
```

Create new migration:
```bash
migrate create -ext sql -dir ./infra/database/migrations -seq migration_name
```

## Schema Overview

### Tables

- **assets**: Trading assets (stocks, crypto, forex, commodities)
- **historical_prices**: OHLCV data for backtesting
- **strategies**: Trading strategy definitions and parameters
- **backtests**: Backtest execution records and results
- **orders**: Order submissions and executions
- **trades**: Completed trades with PnL
- **positions**: Current positions (for live trading)

### Indexes

All tables have appropriate indexes for:
- Foreign key lookups
- Time-based queries
- Status filtering
- Symbol lookups
