# Kasir API Go

A minimalist Cashier API built with Go. This project demonstrates a layered architecture with PostgreSQL database integration.

## Features

- Health check endpoint
- CRUD operations for Products & Categories
- Category existence validation for products
- PostgreSQL database storage
- Swagger API Documentation

## Tech Stack

- **Language:** Go
- **Framework:** Standard Library (`net/http`)
- **Database:** PostgreSQL
- **Configuration:** Viper

## Project Structure

```
kasir-api-go/
├── cmd/kasir-api/          # Application entry point
│   └── main.go
├── internal/
│   ├── config/             # Configuration loading
│   ├── database/           # Database connection
│   ├── handler/            # HTTP handlers
│   ├── models/             # Data models
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic
│   └── utils/              # Utilities (responses, errors)
├── migrations/             # SQL migrations
├── docs/                   # Swagger documentation
└── .env                    # Environment variables
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL

### Configuration

Create a `.env` file in the project root:

```env
APP_NAME=kasir-api
APP_PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/kasir?sslmode=disable
DATABASE_MAX_OPEN_CONNS=25
DATABASE_MAX_IDLE_CONNS=25
```

### Database Setup

1. Create the database:
   ```bash
   createdb kasir
   ```

2. Run migrations:
   ```bash
   psql kasir < migrations/001_create_tables.sql
   ```

### Running Locally

```bash
go run cmd/kasir-api/main.go
```

The server will start on `http://localhost:8080`

## API Documentation

Interactive API documentation (Swagger UI) is available at:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

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

This project is prepared for deployment on [Railway](https://railway.app/). 

1. Connect your GitHub repository to Railway
2. Set environment variables (DATABASE_URL, etc.)
3. Railway will automatically detect and deploy

