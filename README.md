# Kasir API Go

A minimalist Cashier API built with Go. This project demonstrates a layered architecture with PostgreSQL database integration via `pgx`.

## Features

- Health check endpoint
- CRUD operations for Products & Categories
- Category existence validation for products
- PostgreSQL database storage (Supabase compatible)
- Swagger API Documentation with dynamic Host support

## Tech Stack

- **Language:** Go
- **Framework:** Standard Library (`net/http`)
- **Database:** PostgreSQL (Driver: `pgx/v5`)
- **Configuration:** Viper

## Project Structure

```
kasir-api-go/
├── cmd/kasir-api/          # Application entry point
│   └── main.go
├── internal/
│   ├── config/             # Configuration loading
│   ├── database/           # Database connection (PGX)
│   ├── handler/            # HTTP handlers
│   ├── models/             # Data models
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic
│   └── utils/              # Utilities (responses, errors)
├── migrations/             # SQL migrations
├── docs/                   # Swagger documentation
├── railway.json            # Railway deployment config
└── .env                    # Environment variables
```

## Getting Started

### Prerequisites

- Go 1.25 or later
- PostgreSQL (Local or Supabase)

### Configuration

Create a `.env` file in the project root:

```env
APP_NAME=kasir-api
APP_PORT=8080
APP_URL=localhost:8080

# For Local PostgreSQL
DATABASE_URL="postgres://username:password@localhost/kasir?sslmode=disable"

# For Supabase Transaction Pooler (Port 6543)
# default_query_exec_mode=simple_protocol is REQUIRED for transaction pooler
DATABASE_URL="user=postgres.XXXX password=XXXX host=aws-1-ap-south-1.pooler.supabase.com port=6543 dbname=postgres sslmode=require default_query_exec_mode=simple_protocol"

DATABASE_MAX_OPEN_CONNS=1
DATABASE_MAX_IDLE_CONNS=1
```

### Database Setup

1. Create the database:
   ```bash
   createdb kasir
   ```

2. Run migrations:
   ```bash
   psql "YOUR_DATABASE_URL" < migrations/001_create_tables.sql
   ```

### Running Locally

```bash
go run cmd/kasir-api/main.go
```

The server will start on `http://localhost:8080`

## API Documentation

Interactive API documentation (Swagger UI) is available at:
`http://localhost:8080/swagger/index.html`

The host in Swagger is dynamic based on the `APP_URL` environment variable.

### Regenerate Swagger Docs

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/kasir-api/main.go
```

## API Endpoints

### Health Check
`GET /health` - Returns API status

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/products` | Get all products |
| GET | `/api/products/{id}` | Get product by ID |
| POST | `/api/products` | Create product |
| PUT | `/api/products/{id}` | Update product |
| DELETE | `/api/products/{id}` | Delete product |

### Categories
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/categories` | Get all categories |
| GET | `/api/categories/{id}` | Get category by ID |
| POST | `/api/categories` | Create category |
| PUT | `/api/categories/{id}` | Update category |
| DELETE | `/api/categories/{id}` | Delete category |

## Deployment

This project is prepared for deployment on [Railway](https://railway.app/) using the provided `railway.json`.

1. Connect your GitHub repository to Railway.
2. Set environment variables:
   - `APP_URL`: Your Railway domain (e.g., `yourapp.up.railway.app`).
   - `DATABASE_URL`: Your Supabase connection string.
3. Railway will build using the Go builder and the configuration in `railway.json`.

