# Kasir API Go (Sample Project)

A minimalist Cashier API built with Go using only the standard library. This project is created for learning purposes and as a code sample.

## Features

- Health check endpoint
- CRUD operations for Products
- In-memory data storage (restarts reset data)

## Tech Stack

- **Language:** Go
- **Framework:** Standard Library (`net/http`)

## Getting Started

### Prerequisites

- Go 1.25.4 or later

### Running Locally

1. Clone the repository
2. Run the application:
   ```bash
   go run main.go
   ```
3. The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
`GET /health`
- **Response:** `200 OK`
- **Body:** `{"status": "ok", "message": "API Running"}`

### Products

#### Get All Products
`GET /api/products`

#### Get Product by ID
`GET /api/products/{id}`

#### Create Product
`POST /api/products`
- **Body:**
  ```json
  {
    "id": 3,
    "name": "New Product",
    "price": 15000,
    "stock": 50
  }
  ```

#### Update Product
`PUT /api/products/{id}`
- **Body:**
  ```json
  {
    "name": "Updated Product",
    "price": 12000,
    "stock": 40
  }
  ```

#### Delete Product
`DELETE /api/products/{id}`

## Deployment

This project is prepared for deployment on [Railway](https://railway.app/). 

### Step to Deploy:
1. Connect your GitHub repository to Railway.
2. Railway will automatically detect the Go environment and deploy it.
3. Ensure the `PORT` environment variable is handled if necessary (currently hardcoded to `8080`, but Railway usually handles this).

> [!NOTE]
> This is a sample project using in-memory storage. Data will be lost when the server restarts or redeploys.
